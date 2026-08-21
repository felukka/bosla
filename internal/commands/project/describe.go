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
	pname string
}

func NewDescribeCommand() *cobra.Command {
	opts := &describeOptions{}

	describe := &cobra.Command{
		Use:     "describe",
		Args:    cobra.ExactArgs(1),
		Short:   "Describe a project",
		Example: "describe <project_name>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				project, err := describe(ctx, opts, queries)
				if err != nil {
					return err
				}

				output, err = utils.Format(project, output)
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

	return describe
}

func describe(
	ctx context.Context,
	opts *describeOptions,
	q *database.Queries,
) (models.Project, error) {
	projectRecord, err := q.GetProjectByName(ctx, opts.pname)
	if err != nil {
		return models.Project{}, fmt.Errorf(ErrProjectNotFound, opts.pname)
	}

	return models.ToProject(projectRecord), nil
}
