package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/models"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type describeOptions struct {
	name string
}

func NewDescribeCommand(q *database.Queries) *cobra.Command {
	opts := &describeOptions{}

	describe := &cobra.Command{
		Use:     "describe <project>",
		Aliases: []string{"show"},
		Args:    cobra.ExactArgs(1),
		Short:   "Describe a project",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.name = args[0]
		},
		RunE: utils.Wrap(func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
			project, err := describeProject(ctx, opts.name, q)
			if err != nil {
				return err
			}

			formatted, err := utils.Format(project, outputFormat)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), formatted)
			return nil
		}),
	}

	return describe
}

func describeProject(ctx context.Context, projectName string, q *database.Queries) (models.Project, error) {
	projectRecord, err := q.GetProjectByName(ctx, projectName)
	if err != nil {
		return models.Project{}, fmt.Errorf(ErrProjectNotFound, projectName)
	}

	return models.ToProject(projectRecord), nil
}
