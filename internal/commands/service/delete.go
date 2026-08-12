package service

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

func NewDeleteCommand(q *database.Queries) *cobra.Command {
	delete := &cobra.Command{
		Use:     "delete <project_name> <application_name>",
		Args:    cobra.RangeArgs(1, 2),
		Aliases: []string{"remove", "rm"},
		Short:   "Delete an application",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			ctx := cmd.Context()
			pName, aName, err := utils.ParseProjectSlashApplication(args)
			if err != nil {
				return err
			}

			return deleteApplication(ctx, pName, aName, q)
		},
	}

	return delete
}

func deleteApplication(ctx context.Context, pName, aName string, q *database.Queries) error {
	pId, err := q.GetProjectIdByName(ctx, pName)
	if err != nil {
		return err
	}

	if _, err := q.CheckApplicationExistsByName(ctx, database.CheckApplicationExistsByNameParams{
		Name:      aName,
		ProjectID: pId,
	}); err != nil {
		return err
	}

	if _, err = q.DeleteProjectApplicationByName(ctx, database.DeleteProjectApplicationByNameParams{
		Name:      aName,
		ProjectID: pId,
	}); err != nil {
		return fmt.Errorf(ErrServiceDelete, err)
	}

	return err
}
