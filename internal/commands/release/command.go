package release

import (
	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
)

const (
	ErrReleaseExists   = "release %q already exists for project %q"
	ErrReleaseNotFound = "release %q not found for project %q"
)

func NewReleaseCommand(q *database.Queries) *cobra.Command {
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
		NewListCommand(q),
		NewPublishCommand(q),
	)

	return releases
}
