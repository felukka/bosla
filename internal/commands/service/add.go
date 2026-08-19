package service

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type addOptions struct {
	projectName string
	name        string
	status      string
	link        string
	description string
}

func NewAddCommand(q *database.Queries) *cobra.Command {
	opts := &addOptions{}

	add := &cobra.Command{
		Use:     "add <project_name> <service_name>",
		Aliases: []string{"a", "new"},
		Short:   "Add a new service to a project",
		Args:    cobra.RangeArgs(1, 2),
		RunE: utils.Wrap(func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
			var err error
			opts.projectName, opts.name, err = utils.ParseProjectSlashApplication(args)
			if err != nil {
				return err
			}

			return addService(ctx, opts, q)
		}),
	}

	add.Flags().StringVarP(&opts.status, "status", "s", "active", "Service status")
	add.Flags().StringVarP(&opts.link, "link", "l", "", "Service repository URL")
	add.Flags().StringVarP(&opts.description, "description", "d", "", "Service description")

	return add
}

func addService(ctx context.Context, opts *addOptions, q *database.Queries) error {
	projectID, err := q.GetProjectIdByName(ctx, opts.projectName)
	if err != nil {
		return fmt.Errorf("project %q not found", opts.projectName)
	}

	exists, err := q.CheckServiceExistsByName(ctx, database.CheckServiceExistsByNameParams{
		Name:      opts.name,
		ProjectID: projectID,
	})
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf(ErrServiceExists, opts.name)
	}

	now := time.Now().UTC()
	if _, err := q.CreateService(ctx, database.CreateServiceParams{
		ProjectID:   projectID,
		Name:        opts.name,
		Status:      opts.status,
		Description: opts.description,
		RepoUrl:     opts.link,
		CreatedAt:   now,
		UpdatedAt:   now,
	}); err != nil {
		return fmt.Errorf(ErrServiceCreate, err)
	}

	return nil
}
