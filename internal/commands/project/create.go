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
	pname       string
	status      string
	link        string
	description string
}

func NewCreateCommand() *cobra.Command {
	opts := &createOptions{}

	create := &cobra.Command{
		Use:     "create",
		Short:   "Create a new project",
		Args:    cobra.ExactArgs(1),
		Example: "create <project_name>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				return create(ctx, opts, queries)
			},
		),
	}

	create.Flags().StringVarP(&opts.status, "status", "s", "active", "Project status")
	create.Flags().StringVarP(&opts.link, "link", "l", "", "Project link")
	create.Flags().StringVarP(&opts.description, "desc", "d", "", "Project description")

	return create
}

func create(ctx context.Context, opts *createOptions, q *database.Queries) error {
	exists, err := q.CheckProjectExistsByName(ctx, opts.pname)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf(ErrProjectExists, opts.pname)
	}

	now := time.Now().UTC()
	if _, err = q.CreateProject(ctx, database.CreateProjectParams{
		Name:        opts.pname,
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
