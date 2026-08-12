/*
Copyright © 2026 Felukka <felukka.org>
*/
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/commands/config"
	"github.com/mahmoudk1000/bosla/internal/commands/project"
	"github.com/mahmoudk1000/bosla/internal/commands/service"
	"github.com/mahmoudk1000/bosla/internal/database"
)

var bosla = &cobra.Command{
	Use:   "bosla",
	Short: "A serious, well-scoped versioning tool.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.Help(); err != nil {
			return err
		}

		return nil
	},
}

func main() {
	dbURL := os.Getenv("BOSLA_DATABASE_URL")
	if dbURL == "" {
		fmt.Println("Error: BOSLA_DATABASE_URL environment variable is not set")
		fmt.Println(
			"Example: export BOSLA_DATABASE_URL='postgresql://user:pass@localhost:5432/bosla'",
		)

		if err := bosla.Help(); err != nil {
			os.Exit(1)
		}
	}

	if err := database.Init(dbURL); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := database.Close(); err != nil {
			fmt.Printf("Failed to close database connection: %v\n", err)
			os.Exit(1)
		}
	}()

	queries := database.Get()

	bosla.AddCommand(
		project.NewProjectCommand(),
		service.NewServiceCommand(queries),
		config.NewConfigCommand(),
	)

	if err := bosla.Execute(); err != nil {
		os.Exit(1)
	}
}
