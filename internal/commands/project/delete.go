package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type deleteOptions struct {
	name      string
	confirmed bool
}

func NewDeleteCommand(q *database.Queries) *cobra.Command {
	opts := &deleteOptions{}

	deleteCmd := &cobra.Command{
		Use:     "delete <project-name>",
		Aliases: []string{"del", "rm"},
		Short:   "Delete a project",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]
			if !opts.confirmed {
				return fmt.Errorf("deletion requires explicit confirmation: use --yes-i-am-sure flag")
			}
			return nil
		},
		RunE: utils.Wrap(func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
			return deleteProject(ctx, opts, q)
		}),
	}

	deleteCmd.Flags().BoolVar(&opts.confirmed, "yes-i-am-sure", false, "Confirm project deletion")

	return deleteCmd
}

func deleteProject(ctx context.Context, opts *deleteOptions, q *database.Queries) error {
	exists, err := q.CheckProjectExistsByName(ctx, opts.name)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf(ErrProjectNotFound, opts.name)
	}

	return q.DeleteProjectByName(ctx, opts.name)
}
