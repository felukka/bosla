package service

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type registerOptions struct {
	pname       string
	sname       string
	status      string
	link        string
	description string
}

func NewRegisterCommand(q *database.Queries) *cobra.Command {
	opts := &registerOptions{}

	register := &cobra.Command{
		Use:     "register",
		Short:   "Register a new service to a project",
		Args:    cobra.ExactArgs(2),
		Example: "register <project_name> <service_name>",
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
				return registerService(ctx, opts, q)
			},
		),
	}

	register.Flags().StringVarP(&opts.status, "status", "s", "active", "Service status")
	register.Flags().StringVarP(&opts.link, "link", "l", "", "Service repository URL")
	register.Flags().StringVarP(&opts.description, "desc", "d", "", "Service description")

	return register
}

func registerService(ctx context.Context, opts *registerOptions, q *database.Queries) error {
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
	if exists {
		return fmt.Errorf(ErrServiceExists, opts.sname)
	}

	now := time.Now().UTC()
	if _, err := q.CreateService(ctx, database.CreateServiceParams{
		ProjectID:   projectID,
		Name:        opts.sname,
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
