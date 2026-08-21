package project

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type updateOptions struct {
	pname       string
	status      string
	link        string
	description string
}

func NewUpdateCommand() *cobra.Command {
	opts := &updateOptions{}

	update := &cobra.Command{
		Use:     "update",
		Short:   "Update project fields",
		Args:    cobra.ExactArgs(1),
		Example: "update <project_name>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				return update(ctx, opts, queries)
			},
		),
	}

	update.Flags().StringVarP(&opts.status, "status", "s", "", "Project status")
	update.Flags().StringVarP(&opts.link, "link", "l", "", "Project link")
	update.Flags().StringVarP(&opts.description, "description", "d", "", "Project description")

	return update
}

func update(ctx context.Context, opts *updateOptions, q *database.Queries) error {
	exists, err := q.CheckProjectExistsByName(ctx, opts.pname)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf(ErrProjectNotFound, opts.pname)
	}

	current, err := q.GetProjectByName(ctx, opts.pname)
	if err != nil {
		return err
	}

	status := current.Status
	link := current.Link
	description := current.Description

	if opts.status != "" {
		status = opts.status
	}
	if opts.link != "" {
		link = opts.link
	}
	if opts.description != "" {
		description = opts.description
	}

	return q.UpdateProjectByName(ctx, database.UpdateProjectByNameParams{
		Name:        opts.pname,
		Status:      status,
		Link:        link,
		Description: description,
		UpdatedAt:   time.Now().UTC(),
	})
}
