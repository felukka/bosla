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

var bosla = &cobra.Command{
	Use:   "bosla",
	Short: "A serious, well-scoped versioning tool.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.SilenceUsage = true
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.Help(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	bosla.PersistentFlags().StringP("output", "o", "table", "Output format (table|json)")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	dbURL := os.Getenv("BOSLA_DATABASE_URL")
	if dbURL == "" {
		fmt.Println("err: BOSLA_DATABASE_URL env is not set")
		os.Exit(1)
	}

	if err := database.Init(dbURL); err != nil {
		fmt.Printf("failed to connect to database: %v", err)
		os.Exit(1)
	}

	defer func() {
		if err := database.Close(); err != nil {
			fmt.Printf("failed to close database connection: %v\n", err)
			os.Exit(1)
		}
	}()

	queries := database.Get()

	bosla.AddCommand(
		configure.NewConfigureCommand(queries),
		project.NewProjectCommand(queries),
		release.NewReleaseCommand(queries),
		service.NewServiceCommand(queries),
	)

	if err := bosla.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
