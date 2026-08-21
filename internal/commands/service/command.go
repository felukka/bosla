/*
Copyright © 2026 felukka.org
*/
package service

import (
	"github.com/spf13/cobra"
)

const (
	ErrServiceNotFound = "service %q not found"
	ErrServiceExists   = "service %q already exists"
	ErrServiceCreate   = "failed to create service: %w"
	ErrServiceList     = "failed to list services for project %q: %w"
	ErrServiceDelete   = "failed to delete service: %w"
)

func NewServiceCommand() *cobra.Command {
	service := &cobra.Command{
		Use:     "service",
		Aliases: []string{"svc"},
		Short:   "Manage services",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	service.AddCommand(
		NewDeleteCommand(),
		NewListCommand(),
		NewRegisterCommand(),
		NewSearchCommand(),
	)

	return service
}
