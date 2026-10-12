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
	sname string
	key   string
}

func NewDeleteCommand() *cobra.Command {
	opts := &deleteOptions{}

	delete := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm", "del"},
		Short:   "Delete a service metadata key",
		Args:    cobra.ExactArgs(3),
		Example: "delete <project_name> <service_name> <key>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
			opts.sname = args[1]
			opts.key = args[2]
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

	svc, err := q.GetserviceByName(ctx, database.GetserviceByNameParams{
		Name:      opts.sname,
		ProjectID: projectID,
	})
	if err != nil {
		return fmt.Errorf("service %q not found", opts.sname)
	}

	err = q.DeleteServiceMetadata(ctx, database.DeleteServiceMetadataParams{
		ServiceID: svc.ID,
		Key:       opts.key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}

	fmt.Fprintf(out, "Metadata %q deleted from service %q\n", opts.key, opts.sname)
	return nil
}
