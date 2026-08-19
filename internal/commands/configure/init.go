package configure

import (
	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
)

func NewInitCommand(q *database.Queries) *cobra.Command {

	init := &cobra.Command{
		Use:   "init",
		Short: "Initialize database for the CLI tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if err := q.Migrate(ctx); err != nil {
				return err
			}

			return nil
		},
	}

	return init
}
