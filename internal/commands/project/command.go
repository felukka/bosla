package project

import (
	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
)

const (
	ErrProjectNotFound = "project %q not found"
	ErrProjectExists   = "project %q already exists"
	ErrProjectCreate   = "could not create project: %w"
)

func NewProjectCommand(q *database.Queries) *cobra.Command {
	project := &cobra.Command{
		Use:     "project",
		Aliases: []string{"proj", "projects"},
		Short:   "Manage projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	project.AddCommand(
		NewCreateCommand(q),
		NewDeleteCommand(q),
		NewDescribeCommand(q),
		NewGetCommand(q),
		NewSearchCommand(q),
		NewUpdateCommand(q),
	)

	return project
}
