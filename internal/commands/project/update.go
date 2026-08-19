package project

import (
	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
)

func NewUpdateCommand(q *database.Queries) *cobra.Command {
	update := &cobra.Command{
		Use:   "update <project>",
		Short: "update projects",
		Args:  cobra.ExactArgs(1),
	}

	return update
}
