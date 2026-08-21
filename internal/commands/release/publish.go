package release

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type publishOptions struct {
	pname       string
	version     string
	description string
	status      string
	gitHash     string
	services    []string
}

func NewPublishCommand() *cobra.Command {
	opts := &publishOptions{}

	publish := &cobra.Command{
		Use:     "release",
		Short:   "Release a project release and attach service versions",
		Args:    cobra.ExactArgs(2),
		Example: "release <project_name> <version> [--service <service_name>]",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
			opts.version = args[1]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				return publish(ctx, opts, queries)
			},
		),
	}

	publish.Flags().StringVarP(&opts.description, "description", "d", "", "Release description")
	publish.Flags().StringVarP(&opts.status, "status", "s", "released", "Service version status")
	publish.Flags().StringVar(&opts.gitHash, "git-hash", "", "Service version git hash")
	publish.Flags().
		StringArrayVar(&opts.services, "service", []string{}, "Attach service version as service=version (repeatable)")

	return publish
}

func publish(ctx context.Context, opts *publishOptions, q *database.Queries) error {
	projectID, err := q.GetProjectIdByName(ctx, opts.pname)
	if err != nil {
		return fmt.Errorf("project %q not found", opts.pname)
	}

	exists, err := q.CheckProjectVersionExists(ctx, database.CheckProjectVersionExistsParams{
		ProjectID: projectID,
		Version:   opts.version,
	})
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf(ErrReleaseExists, opts.version, opts.pname)
	}

	now := time.Now().UTC()
	projectVersion, err := q.CreateProjectVersion(ctx, database.CreateProjectVersionParams{
		ProjectID:   projectID,
		Version:     opts.version,
		Description: opts.description,
		CreatedAt:   now,
	})
	if err != nil {
		return err
	}

	serviceVersionOverrides, err := parseServiceVersionArgs(opts.services)
	if err != nil {
		return err
	}

	services, err := q.GetServiceByProjectName(ctx, projectID)
	if err != nil {
		return err
	}

	if len(serviceVersionOverrides) > 0 {
		allowed := make(map[string]struct{}, len(services))
		for _, svc := range services {
			allowed[svc.Name] = struct{}{}
		}
		for sname := range serviceVersionOverrides {
			if _, ok := allowed[sname]; !ok {
				return fmt.Errorf(
					"service %q does not belong to project %q",
					sname,
					opts.pname,
				)
			}
		}
	}

	for _, svc := range services {
		svcVersion := opts.version
		if len(serviceVersionOverrides) > 0 {
			v, ok := serviceVersionOverrides[svc.Name]
			if !ok {
				continue
			}
			svcVersion = v
		}

		serviceVersionID, err := ensureServiceVersion(
			ctx,
			q,
			svc.ID,
			svcVersion,
			opts.status,
			opts.gitHash,
			opts.description,
			now,
		)
		if err != nil {
			return err
		}

		if err := q.UpsertProjectVersionService(ctx, database.UpsertProjectVersionServiceParams{
			ProjectVersionID: projectVersion.ID,
			ServiceID:        svc.ID,
			ServiceVersionID: serviceVersionID,
		}); err != nil {
			return err
		}
	}

	return nil
}

func ensureServiceVersion(
	ctx context.Context,
	q *database.Queries,
	serviceID int32,
	version string,
	status string,
	gitHash string,
	description string,
	createdAt time.Time,
) (int32, error) {
	exists, err := q.CheckServiceVersionExists(ctx, database.CheckServiceVersionExistsParams{
		ServiceID: serviceID,
		Version:   version,
	})
	if err != nil {
		return 0, err
	}

	if exists {
		sv, err := q.GetServiceVersionByServiceAndVersion(
			ctx,
			database.GetServiceVersionByServiceAndVersionParams{
				ServiceID: serviceID,
				Version:   version,
			},
		)
		if err != nil {
			return 0, err
		}
		return sv.ID, nil
	}

	sv, err := q.CreateServiceVersion(ctx, database.CreateServiceVersionParams{
		ServiceID:   serviceID,
		Version:     version,
		Status:      status,
		GitHash:     gitHash,
		Description: description,
		CreatedAt:   createdAt,
	})
	if err != nil {
		return 0, err
	}

	return sv.ID, nil
}

func parseServiceVersionArgs(values []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, value := range values {
		parts := strings.SplitN(value, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid --service value %q: expected service=version", value)
		}

		serviceName := strings.TrimSpace(parts[0])
		version := strings.TrimSpace(parts[1])
		if serviceName == "" || version == "" {
			return nil, fmt.Errorf("invalid --service value %q: expected service=version", value)
		}

		result[serviceName] = version
	}
	return result, nil
}
