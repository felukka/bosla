/*
Copyright © 2026 (mahmoudk1000) <mahmoudk1000@gmail.com>
*/
package service

import (
	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
)

const (
	ErrServiceNotFound = "service %q not found"
	ErrServiceExists   = "service %q already exists"
	ErrServiceCreate   = "failed to create service: %w"
	ErrServiceList     = "failed to list services for project %q: %w"
	ErrServiceDelete   = "failed to delete service: %w"
)

func NewServiceCommand(q *database.Queries) *cobra.Command {
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
		NewDeleteCommand(q),
		NewListCommand(q),
		NewRegisterCommand(q),
		NewSearchCommand(q),
	)

	return service
}
