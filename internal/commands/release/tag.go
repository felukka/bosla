package version

import "github.com/spf13/cobra"

func NewCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "create <version>",
		Aliases: []string{"c", "new"},
		Short:   "Create a new version",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}
