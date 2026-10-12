package service

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
	sname       string
	status      string
	link        string
	description string
}

func NewUpdateCommand() *cobra.Command {
	opts := &updateOptions{}

	update := &cobra.Command{
		Use:     "update",
		Args:    cobra.ExactArgs(2),
		Short:   "Update service fields",
		Example: "update <project_name> <service_name>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
			opts.sname = args[1]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				return update(ctx, opts, queries)
			},
		),
	}

	update.Flags().StringVarP(&opts.status, "status", "s", "", "Service status")
	update.Flags().StringVarP(&opts.link, "link", "l", "", "Service repository URL")
	update.Flags().StringVarP(&opts.description, "desc", "d", "", "Service description")

	return update
}

func update(ctx context.Context, opts *updateOptions, q *database.Queries) error {
	projectID, err := q.GetProjectIdByName(ctx, opts.pname)
	if err != nil {
		return fmt.Errorf("project %q not found", opts.pname)
	}

	svc, err := q.GetserviceByName(ctx, database.GetserviceByNameParams{
		Name:      opts.sname,
		ProjectID: projectID,
	})
	if err != nil {
		return fmt.Errorf(ErrServiceNotFound, opts.sname)
	}

	status := svc.Status
	link := svc.RepoUrl
	description := svc.Description

	if opts.status != "" {
		status = opts.status
	}
	if opts.link != "" {
		link = opts.link
	}
	if opts.description != "" {
		description = opts.description
	}

	now := time.Now().UTC()
	return q.UpdateServiceByName(ctx, database.UpdateServiceByNameParams{
		ID:          svc.ID,
		Status:      status,
		RepoUrl:     link,
		Description: description,
		UpdatedAt:   now,
	})
}
