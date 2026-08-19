package project

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type createOptions struct {
	name        string
	status      string
	link        string
	description string
}

func NewCreateCommand(q *database.Queries) *cobra.Command {
	opts := &createOptions{}

	create := &cobra.Command{
		Use:     "create <name>",
		Aliases: []string{"c", "new"},
		Short:   "add a new application to the project",
		Args:    cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.name = args[0]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
				return createProject(ctx, opts, q)
			},
		),
	}

	create.Flags().StringVarP(&opts.status, "status", "s", "active", "Project status")
	create.Flags().StringVarP(&opts.link, "link", "l", "", "Project link")
	create.Flags().StringVarP(&opts.description, "description", "d", "", "Project description")

	return create
}

func createProject(
	ctx context.Context,
	opts *createOptions,
	q *database.Queries,
) error {
	exists, err := q.CheckProjectExistsByName(ctx, opts.name)
	if err != nil {
		return fmt.Errorf(ErrProjectNotFound, err)
	}
	if exists {
		return fmt.Errorf(ErrProjectNotFound, opts.name)
	}

	now := time.Now().UTC()
	if _, err = q.CreateProject(ctx, database.CreateProjectParams{
		Name:        opts.name,
		Status:      opts.status,
		Link:        opts.link,
		Description: opts.description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}); err != nil {
		return fmt.Errorf(ErrProjectCreate, err)
	}

	return nil
}
