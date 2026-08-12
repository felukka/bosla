/*
Copyright © 2026 mahmoudk1000 <mahmoudk1000@gmail.com>
*/
package project

import (
	"github.com/spf13/cobra"
)

const (
	ErrProjectNotFound = "failed to find project: %w"
)

func NewProjectCommand() *cobra.Command {
	project := &cobra.Command{
		Use:     "project create|delete|list|metadata|show|status",
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
		NewCreateCommand(),
		NewDeleteCommand(),
		NewListCommand(),
		NewMetadataCommand(),
		NewShowCommand(),
		NewStatusCommand(),
	)

	return project
}
