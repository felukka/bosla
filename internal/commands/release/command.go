package version

import (
	"github.com/spf13/cobra"
)

const (
	failedToParseVersionErr = "failed to parse version: %w"
)

func NewVersionCommand() *cobra.Command {
	versions := &cobra.Command{
		Use:     "version create|delete|list|remove|show",
		Short:   "Manage versions",
		Aliases: []string{"ver"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}

			return nil
		},
	}

	versions.AddCommand()

	return versions
}
