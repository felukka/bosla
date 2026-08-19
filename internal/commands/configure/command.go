package configure

import (
	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
)

func NewConfigureCommand(q *database.Queries) *cobra.Command {
	configure := &cobra.Command{
		Use:     "configure",
		Aliases: []string{"config"},
		Short:   "manage configuration for bosla",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}

			return nil
		},
	}

	configure.AddCommand(NewInitCommand(q))

	return configure
}
