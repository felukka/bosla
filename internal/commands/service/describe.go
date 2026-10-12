package service

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/models"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type describeOptions struct {
	pname string
	sname string
}

func NewDescribeCommand() *cobra.Command {
	opts := &describeOptions{}

	describe := &cobra.Command{
		Use:     "describe",
		Args:    cobra.ExactArgs(2),
		Short:   "Describe a service",
		Example: "describe <project_name> <service_name>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
			opts.sname = args[1]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				svc, err := describe(ctx, opts, queries)
				if err != nil {
					return err
				}

				formatted, err := utils.Format(svc, output)
				if err != nil {
					return err
				}

				if _, err := fmt.Fprintln(cmd.OutOrStdout(), formatted); err != nil {
					return err
				}

				return nil
			},
		),
	}

	return describe
}

func describe(
	ctx context.Context,
	opts *describeOptions,
	q *database.Queries,
) (*models.ServiceOutput, error) {
	projectID, err := q.GetProjectIdByName(ctx, opts.pname)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", opts.pname)
	}

	svc, err := q.GetserviceByName(ctx, database.GetserviceByNameParams{
		Name:      opts.sname,
		ProjectID: projectID,
	})
	if err != nil {
		return nil, fmt.Errorf(ErrServiceNotFound, opts.sname)
	}

	metadata, err := q.GetServiceMetadata(ctx, svc.ID)
	if err != nil {
		return nil, err
	}

	meta := make(map[string]string, len(metadata))
	for _, m := range metadata {
		meta[m.Key] = m.Value
	}

	return &models.ServiceOutput{
		Name:        svc.Name,
		Status:      svc.Status,
		RepoURL:     svc.RepoUrl,
		Description: svc.Description,
		CreatedAt:   svc.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   svc.UpdatedAt.Format(time.RFC3339),
		Metadata:    meta,
	}, nil
}
