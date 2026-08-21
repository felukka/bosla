package service

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/models"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type listOptions struct {
	pname string
}

func NewListCommand() *cobra.Command {
	opts := &listOptions{}

	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List services for a project",
		Args:    cobra.ExactArgs(1),
		Example: "list <project_name>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				services, err := list(ctx, opts, queries)
				if err != nil {
					return err
				}

				output, err = utils.Format(services, output)
				if err != nil {
					return err
				}

				if _, err := fmt.Fprintln(cmd.OutOrStdout(), output); err != nil {
					return err
				}

				return nil
			},
		),
	}
}

func list(
	ctx context.Context,
	opts *listOptions,
	q *database.Queries,
) ([]models.Service, error) {
	projectID, err := q.GetProjectIdByName(ctx, opts.pname)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", opts.pname)
	}

	services, err := q.GetServiceByProjectName(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf(ErrServiceList, opts.pname, err)
	}

	return models.ToServices(services), nil
}
