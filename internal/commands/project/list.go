package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/models"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

func NewListCommand(q *database.Queries) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List projects",
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, format string) error {
				projects, err := listProjects(ctx, q)
				if err != nil {
					return err
				}

				formatted, err := utils.Format(projects, format)
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
}

func listProjects(ctx context.Context, q *database.Queries) ([]models.Project, error) {
	p, err := q.GetProjects(ctx)
	if err != nil {
		return nil, err
	}

	return models.ToProjects(p), nil
}
