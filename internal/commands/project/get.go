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
		Aliases: []string{"ls"},
		Short:   "Get projects",
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {

				project, err := getProjects(ctx, q)
				if err != nil {
					return err
				}

				format, err := utils.Format(project, outputFormat)
				if err != nil {
					return err
				}

				fmt.Fprintln(cmd.OutOrStdout(), format)
				return nil
			},
		),
	}
}

func getProjects(
	ctx context.Context,
	q *database.Queries,
) ([]models.Project, error) {

	p, err := q.GetProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf(ErrProjectNotFound, err)
	}

	return models.ToProjects(p), nil
}
