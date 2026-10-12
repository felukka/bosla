package metadata

import (
	"github.com/spf13/cobra"
)

func NewMetadataCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "metadata",
		Aliases: []string{"meta", "m"},
		Short:   "Manage service metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		NewListCommand(),
		NewGetCommand(),
		NewSetCommand(),
		NewDeleteCommand(),
	)

	return cmd
}
