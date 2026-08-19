package project

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
)

type deleteOptions struct {
	name      string
	confirmed bool
}

func NewDeleteCommand(q *database.Queries) *cobra.Command {
	opts := &deleteOptions{}

	delete := &cobra.Command{
		Use:     "delete <project-name>",
		Aliases: []string{"del", "rm"},
		Short:   "Delete a project",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]

			if !opts.confirmed {
				return fmt.Errorf(
					"deletion requires explicit confirmation: use --yes-i-am-sure flag",
				)
			}

			return nil
		},
	}

	flags := delete.Flags()
	flags.BoolVar(&opts.confirmed, "yes-i-am-sure", false, "Confirm project deletion")

	delete.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return runDelete(cmd.Context(), opts, q)
	}

	return delete
}

func runDelete(ctx context.Context, opts *deleteOptions, q *database.Queries) error {
	exists, err := q.CheckProjectExistsByName(ctx, opts.name)

	if err != nil {
		return fmt.Errorf(ErrProjectNotFound, opts.name, err)
	}
	if !exists {
		return fmt.Errorf(ErrProjectNotFound, opts.name)
	}

	return q.DeleteProjectByName(ctx, opts.name)
}
