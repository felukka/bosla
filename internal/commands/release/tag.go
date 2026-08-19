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

type createOptions struct {
	projectName string
	version     string
	description string
	status      string
	gitHash     string
	services    []string
}

func NewCreateCommand(q *database.Queries) *cobra.Command {
	opts := &createOptions{}

	create := &cobra.Command{
		Use:     "create <project> <release-version>",
		Aliases: []string{"new"},
		Short:   "Create a project release and attach service versions",
		Args:    cobra.ExactArgs(2),
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.projectName = args[0]
			opts.version = args[1]
		},
		RunE: utils.Wrap(func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
			return createRelease(ctx, opts, q)
		}),
	}

	create.Flags().StringVarP(&opts.description, "description", "d", "", "Release description")
	create.Flags().StringVarP(&opts.status, "status", "s", "released", "Service version status")
	create.Flags().StringVar(&opts.gitHash, "git-hash", "", "Service version git hash")
	create.Flags().StringArrayVar(&opts.services, "service", []string{}, "Attach service version as service=version (repeatable)")

	return create
}

func createRelease(ctx context.Context, opts *createOptions, q *database.Queries) error {
	projectID, err := q.GetProjectIdByName(ctx, opts.projectName)
	if err != nil {
		return fmt.Errorf("project %q not found", opts.projectName)
	}

	exists, err := q.CheckProjectVersionExists(ctx, database.CheckProjectVersionExistsParams{
		ProjectID: projectID,
		Version:   opts.version,
	})
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf(ErrReleaseExists, opts.version, opts.projectName)
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
		for serviceName := range serviceVersionOverrides {
			if _, ok := allowed[serviceName]; !ok {
				return fmt.Errorf("service %q does not belong to project %q", serviceName, opts.projectName)
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

		serviceVersionID, err := ensureServiceVersion(ctx, q, svc.ID, svcVersion, opts.status, opts.gitHash, opts.description, now)
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
		sv, err := q.GetServiceVersionByServiceAndVersion(ctx, database.GetServiceVersionByServiceAndVersionParams{
			ServiceID: serviceID,
			Version:   version,
		})
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
