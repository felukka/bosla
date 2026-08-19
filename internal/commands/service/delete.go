package service

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

func NewDeleteCommand(q *database.Queries) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:     "delete <project_name> <service_name>",
		Args:    cobra.RangeArgs(1, 2),
		Aliases: []string{"remove", "rm"},
		Short:   "Delete a service",
		RunE: utils.Wrap(func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
			pName, sName, err := utils.ParseProjectSlashApplication(args)
			if err != nil {
				return err
			}

			return deleteService(ctx, pName, sName, q)
		}),
	}

	return deleteCmd
}

func deleteService(ctx context.Context, projectName, serviceName string, q *database.Queries) error {
	projectID, err := q.GetProjectIdByName(ctx, projectName)
	if err != nil {
		return fmt.Errorf("project %q not found", projectName)
	}

	exists, err := q.CheckServiceExistsByName(ctx, database.CheckServiceExistsByNameParams{
		Name:      serviceName,
		ProjectID: projectID,
	})
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf(ErrServiceNotFound, serviceName)
	}

	if _, err = q.DeleteProjectServiceByName(ctx, database.DeleteProjectServiceByNameParams{
		Name:      serviceName,
		ProjectID: projectID,
	}); err != nil {
		return fmt.Errorf(ErrServiceDelete, err)
	}

	return nil
}
