package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/models"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

func NewGetCommand(q *database.Queries) *cobra.Command {
	return &cobra.Command{
		Use:     "get",
		Aliases: []string{"ls", "list"},
		Short:   "List projects",
		RunE: utils.Wrap(func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
			projects, err := getProjects(ctx, q)
			if err != nil {
				return err
			}

			formatted, err := utils.Format(projects, outputFormat)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), formatted)
			return nil
		}),
	}
}

func getProjects(ctx context.Context, q *database.Queries) ([]models.Project, error) {
	p, err := q.GetProjects(ctx)
	if err != nil {
		return nil, err
	}

	return models.ToProjects(p), nil
}
