package project

import (
	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/commands/project/metadata"
)

const (
	ErrProjectNotFound = "project %q not found"
	ErrProjectExists   = "project %q already exists"
	ErrProjectCreate   = "could not create project: %w"
)

func NewProjectCommand() *cobra.Command {
	project := &cobra.Command{
		Use:     "project",
		Aliases: []string{"prj"},
		Short:   "Manage projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	project.AddCommand(
		NewCreateCommand(),
		NewDeleteCommand(),
		NewDescribeCommand(),
		NewListCommand(),
		NewSearchCommand(),
		NewUpdateCommand(),
		metadata.NewMetadataCommand(),
	)

	return project
}
