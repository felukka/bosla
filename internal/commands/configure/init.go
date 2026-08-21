package configure

import (
	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
)

func NewInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize database for the CLI tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			q := database.Get()

			if err := q.Migrate(ctx); err != nil {
				return err
			}

			return nil
		},
	}
}
