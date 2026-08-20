package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type deleteOptions struct {
	pname   string
	confirm bool
}

func NewDeleteCommand(q *database.Queries) *cobra.Command {
	opts := &deleteOptions{}

	delete := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Short:   "Delete a project",
		Args:    cobra.ExactArgs(1),
		Example: "delete <project_name>",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.pname = args[0]
			if !opts.confirm {
				return fmt.Errorf(
					"deletion requires explicit confirmation: use --yes-i-am-sure flag",
				)
			}

			return nil
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
				return deleteProject(ctx, opts, q)
			},
		),
	}

	delete.Flags().BoolVar(&opts.confirm, "yes-i-am-sure", false, "confirm project deletion")

	return delete
}

func deleteProject(ctx context.Context, opts *deleteOptions, q *database.Queries) error {
	exists, err := q.CheckProjectExistsByName(ctx, opts.pname)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf(ErrProjectNotFound, opts.pname)
	}

	return q.DeleteProjectByName(ctx, opts.pname)
}
