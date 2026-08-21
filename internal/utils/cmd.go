package utils

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
)

type handler func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error

func Wrap(handler handler) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		queries := database.Get()

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		switch output {
		case "json", "table":
		default:
			return fmt.Errorf("invalid --output: %s (want json|table)", output)
		}

		return handler(ctx, cmd, args, queries, output)
	}
}
