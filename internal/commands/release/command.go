package release

import (
	"github.com/spf13/cobra"
)

const (
	ErrReleaseExists   = "release %q already exists for project %q"
	ErrReleaseNotFound = "release %q not found for project %q"
)

func NewReleaseCommand() *cobra.Command {
	releases := &cobra.Command{
		Use:   "release",
		Short: "Manage project releases",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	releases.AddCommand(
		NewListCommand(),
		NewPublishCommand(),
	)

	return releases
}
