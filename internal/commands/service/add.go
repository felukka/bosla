package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/sqlc-dev/pqtype"

	"github.com/mahmoudk1000/bosla/internal/commands/project"
	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type addOptions struct {
	projectName  string
	app          string
	link         string
	description  string
	metadata     []string
	metadataJSON pqtype.NullRawMessage
}

func NewAddCommand(q *database.Queries) *cobra.Command {
	opts := &addOptions{}

	add := &cobra.Command{
		Use:     "add <project_name> <application_name>",
		Aliases: []string{"a", "new"},
		Short:   "add a new application to the project",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx := cmd.Context()
			var err error

			opts.projectName, opts.app, err = utils.ParseProjectSlashApplication(args)
			if err != nil {
				return err
			}

			metadataMap, err := utils.ParseMetadata(opts.metadata)
			if err != nil {
				return fmt.Errorf(ErrServiceParseMetadata, err)
			}

			opts.metadataJSON, err = utils.MetadataToJSON(metadataMap)
			if err != nil {
				return err
			}

			return addApplication(ctx, opts, q)
		},
	}

	add.Flags().StringP("output", "o", "table", "Output format (table|json)")
	add.Flags().StringVarP(&opts.link, "link", "l", "", "application's link")
	add.Flags().StringVarP(&opts.description, "description", "d", "", "application's description")
	add.Flags().
		StringArrayVarP(&opts.metadata, "metadata", "m", []string{}, "Metadata key=value pairs")

	return add
}

func addApplication(
	ctx context.Context,
	opts *addOptions,
	q *database.Queries,
) error {
	pID, err := q.GetProjectIdByName(ctx, opts.projectName)
	if err != nil {
		return fmt.Errorf(project.ErrProjectNotFound, opts.projectName, err)
	}

	exists, err := q.CheckApplicationExistsByName(ctx, database.CheckApplicationExistsByNameParams{
		Name:      opts.app,
		ProjectID: pID,
	})
	if err != nil {
		return fmt.Errorf(ErrServiceNotFound, err)
	}
	if exists {
		return fmt.Errorf(ErrServiceNotFound, opts.app)
	}

	if _, err := q.CreateApplication(ctx, database.CreateApplicationParams{
		Name:      opts.app,
		ProjectID: pID,
		RepoUrl: sql.NullString{
			String: opts.link,
			Valid:  opts.link != "",
		},
		Description: sql.NullString{
			String: opts.description,
			Valid:  opts.description != "",
		},
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		return fmt.Errorf(ErrServiceCreate, err)
	}

	return nil
}
