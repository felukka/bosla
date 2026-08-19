package utils

import (
	"context"

	"github.com/spf13/cobra"
)

type handler func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error

func Wrap(handler handler) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		outputFormat, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		return handler(ctx, cmd, args, outputFormat)
	}
}
