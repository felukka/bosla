package service

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type deleteOptions struct {
	pname string
	sname string
}

func NewDeleteCommand() *cobra.Command {
	opts := &deleteOptions{}

	delete := &cobra.Command{
		Use:     "delete",
		Args:    cobra.ExactArgs(2),
		Aliases: []string{"rm"},
		Short:   "Delete a service",
		Example: "delete <project_name> <service_name>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
			opts.sname = args[1]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				return delete(ctx, opts, queries)
			},
		),
	}

	return delete
}

func delete(
	ctx context.Context,
	opts *deleteOptions,
	q *database.Queries,
) error {
	projectID, err := q.GetProjectIdByName(ctx, opts.pname)
	if err != nil {
		return fmt.Errorf("project %q not found", opts.pname)
	}

	exists, err := q.CheckServiceExistsByName(ctx, database.CheckServiceExistsByNameParams{
		Name:      opts.sname,
		ProjectID: projectID,
	})
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf(ErrServiceNotFound, opts.sname)
	}

	if _, err = q.DeleteProjectServiceByName(ctx, database.DeleteProjectServiceByNameParams{
		Name:      opts.sname,
		ProjectID: projectID,
	}); err != nil {
		return fmt.Errorf(ErrServiceDelete, err)
	}

	return nil
}
