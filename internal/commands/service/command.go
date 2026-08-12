/*
Copyright © 2026 (mahmoudk1000) <mahmoudk1000@gmail.com>
*/
package service

import (
	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
)

const (
	ErrServiceNotFound      = "failed to find service: %w"
	ErrServiceCreate        = "failed to create service"
	ErrServiceList          = "failed to list application for project %q: %w"
	ErrServiceDelete        = "failed to delete service"
	ErrServiceParseMetadata = "failed to parse service metadata"
)

func NewServiceCommand(q *database.Queries) *cobra.Command {
	service := &cobra.Command{
		Use:     "application add|remove|list|show",
		Aliases: []string{"app", "applications"},
		Short:   "Manage applications",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	service.AddCommand(
		NewAddCommand(q),
		NewDeleteCommand(q),
		NewListCommand(q),
	)

	return service
}
