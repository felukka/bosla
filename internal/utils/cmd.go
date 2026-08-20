package utils

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type handler func(ctx context.Context, cmd *cobra.Command, args []string, Format string) error

func Wrap(handler handler) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		format, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		switch format {
		case "json", "table":
		default:
			return fmt.Errorf("invalid --output: %s (want json|table)", format)
		}

		return handler(ctx, cmd, args, format)
	}
}
