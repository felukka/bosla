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
	name        string
	status      string
	link        string
	description string
}

func NewUpdateCommand(q *database.Queries) *cobra.Command {
	opts := &updateOptions{}

	update := &cobra.Command{
		Use:   "update <project>",
		Short: "Update project fields",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.name = args[0]
		},
		RunE: utils.Wrap(func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
			return updateProject(ctx, opts, q)
		}),
	}

	update.Flags().StringVarP(&opts.status, "status", "s", "", "Project status")
	update.Flags().StringVarP(&opts.link, "link", "l", "", "Project link")
	update.Flags().StringVarP(&opts.description, "description", "d", "", "Project description")

	return update
}

func updateProject(ctx context.Context, opts *updateOptions, q *database.Queries) error {
	exists, err := q.CheckProjectExistsByName(ctx, opts.name)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf(ErrProjectNotFound, opts.name)
	}

	current, err := q.GetProjectByName(ctx, opts.name)
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
		Name:        opts.name,
		Status:      status,
		Link:        link,
		Description: description,
		UpdatedAt:   time.Now().UTC(),
	})
}
