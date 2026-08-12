package version

import (
	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
)

type createOptions struct {
	version string
	status  string
	desc    string
	gitHash string
}

func NewCreateCommand() *cobra.Command {
	opts := &createOptions{}
	var queries *database.Queries

	create := &cobra.Command{
		Use:     "create <version>",
		Aliases: []string{"c", "new"},
		Short:   "Create a new version",
		Args:    cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = database.Get()
		},
	}

	return create
}
