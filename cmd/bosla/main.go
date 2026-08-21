/*
Copyright © 2026 Felukka <felukka.org>
*/
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/commands/configure"
	"github.com/mahmoudk1000/bosla/internal/commands/project"
	"github.com/mahmoudk1000/bosla/internal/commands/release"
	"github.com/mahmoudk1000/bosla/internal/commands/service"
	"github.com/mahmoudk1000/bosla/internal/database"
)

var (
	ErrDatabaseClosing    = "failed to close database connection: %w"
	ErrDatabaseConnection = "failed to connect to database: %w"
)

var bosla = &cobra.Command{
	Use:   "bosla",
	Short: "A serious, well-scoped versioning tool.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.SilenceUsage = true
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		url := os.Getenv("BOSLA_DATABASE_URL")

		if url == "" {
			cfg, err := configure.LoadConfig()
			if err != nil {
				return fmt.Errorf(
					"database connection is not set.\nSet BOSLA_DATABASE_URL or run `bosla configure set --help`",
				)
			}

			port := cfg.Port
			if port == 0 {
				port = 5432
			}

			sslMode := "require"
			if cfg.SkipSSL {
				sslMode = "disable"
			}

			url = fmt.Sprintf("postgres://%s:%s@%s:%d/bosla?sslmode=%s",
				cfg.Username, cfg.Password, cfg.URL, port, sslMode)
		}

		if err := database.Init(url); err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if err := database.Close(); err != nil {
			return fmt.Errorf(ErrDatabaseClosing, err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	bosla.PersistentFlags().StringP("output", "o", "table", "Output format (table|json)")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	bosla.AddCommand(
		configure.NewConfigureCommand(),
		project.NewProjectCommand(),
		release.NewReleaseCommand(),
		service.NewServiceCommand(),
	)

	if err := bosla.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
