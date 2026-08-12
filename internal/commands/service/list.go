package service

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/commands/project"
	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/models"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

func NewListCommand(q *database.Queries) *cobra.Command {
	list := &cobra.Command{
		Use:     "list [project-name]",
		Aliases: []string{"ls"},
		Short:   "List all applications of a project",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			ctx := cmd.Context()
			projectName := args[0]

			outputFormat, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}

			applications, err := listApplications(ctx, projectName, q)
			if err != nil {
				return err
			}

			fmtA, err := utils.Format(applications, outputFormat)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), fmtA)

			return nil
		},
	}

	list.Flags().StringP("output", "o", "table", "Output format (table|json)")

	return list
}

func listApplications(
	ctx context.Context,
	pName string,
	q *database.Queries,
) ([]models.Application, error) {
	pId, err := q.GetProjectIdByName(ctx, pName)
	if err != nil {
		return nil, fmt.Errorf(project.ErrProjectNotFound, pName, err)
	}

	ps, err := q.ListAllProjectApplications(ctx, pId)
	if err != nil {
		return nil, fmt.Errorf(ErrServiceList, pName, err)
	}

	return models.ToApplications(ps), nil
}
