package metadata

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type setOptions struct {
	pname string
	key   string
	value string
}

func NewSetCommand() *cobra.Command {
	opts := &setOptions{}

	set := &cobra.Command{
		Use:     "set",
		Short:   "Set a project metadata key-value pair",
		Args:    cobra.ExactArgs(3),
		Example: "set <project_name> <key> <value>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
			opts.key = args[1]
			opts.value = args[2]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				return set(ctx, opts, queries, cmd.OutOrStdout())
			},
		),
	}

	return set
}

func set(ctx context.Context, opts *setOptions, q *database.Queries, out io.Writer) error {
	projectID, err := q.GetProjectIdByName(ctx, opts.pname)
	if err != nil {
		return fmt.Errorf("project %q not found", opts.pname)
	}

	err = q.UpsertProjectMetadata(ctx, database.UpsertProjectMetadataParams{
		ProjectID: projectID,
		Key:       opts.key,
		Value:     opts.value,
	})
	if err != nil {
		return fmt.Errorf("failed to set metadata: %w", err)
	}

	fmt.Fprintf(out, "Metadata %q set for project %q\n", opts.key, opts.pname)
	return nil
}
