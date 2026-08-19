package service

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/models"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type getOptions struct {
	projectName string
}

func NewListCommand(q *database.Queries) *cobra.Command {
	opts := &getOptions{}

	return &cobra.Command{
		Use:     "get <project>",
		Aliases: []string{"ls", "list"},
		Short:   "List services for a project",
		Args:    cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.projectName = args[0]
		},
		RunE: utils.Wrap(func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
			services, err := getServices(ctx, opts, q)
			if err != nil {
				return err
			}

			formatted, err := utils.Format(services, outputFormat)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), formatted)
			return nil
		}),
	}
}

func getServices(ctx context.Context, opts *getOptions, q *database.Queries) ([]models.Service, error) {
	projectID, err := q.GetProjectIdByName(ctx, opts.projectName)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", opts.projectName)
	}

	services, err := q.GetServiceByProjectName(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf(ErrServiceList, opts.projectName, err)
	}

	return models.ToServices(services), nil
}
