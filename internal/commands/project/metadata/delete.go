package metadata

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type deleteOptions struct {
	pname string
	key   string
}

func NewDeleteCommand() *cobra.Command {
	opts := &deleteOptions{}

	delete := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm", "del"},
		Short:   "Delete a project metadata key",
		Args:    cobra.ExactArgs(2),
		Example: "delete <project_name> <key>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
			opts.key = args[1]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				return delete(ctx, opts, queries, cmd.OutOrStdout())
			},
		),
	}

	return delete
}

func delete(ctx context.Context, opts *deleteOptions, q *database.Queries, out io.Writer) error {
	projectID, err := q.GetProjectIdByName(ctx, opts.pname)
	if err != nil {
		return fmt.Errorf("project %q not found", opts.pname)
	}

	err = q.DeleteProjectMetadata(ctx, database.DeleteProjectMetadataParams{
		ProjectID: projectID,
		Key:       opts.key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}

	fmt.Fprintf(out, "Metadata %q deleted from project %q\n", opts.key, opts.pname)
	return nil
}
