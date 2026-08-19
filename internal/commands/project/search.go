package project

import (
	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
)

func NewSearchCommand(q *database.Queries) *cobra.Command {

	search := &cobra.Command{
		Use:   "search <regex>",
		Short: "search projects",
		Args:  cobra.ExactArgs(1),
	}

	return search
}
