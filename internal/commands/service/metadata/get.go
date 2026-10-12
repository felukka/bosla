package metadata

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type getOptions struct {
	pname string
	sname string
	key   string
}

func NewGetCommand() *cobra.Command {
	opts := &getOptions{}

	get := &cobra.Command{
		Use:     "get",
		Short:   "Get a service metadata value by key",
		Args:    cobra.ExactArgs(3),
		Example: "get <project_name> <service_name> <key>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
			opts.sname = args[1]
			opts.key = args[2]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				return get(ctx, opts, queries, output, cmd.OutOrStdout())
			},
		),
	}

	return get
}

func get(ctx context.Context, opts *getOptions, q *database.Queries, output string, out io.Writer) error {
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

	m, err := q.GetServiceMetadataByKey(ctx, database.GetServiceMetadataByKeyParams{
		ServiceID: svc.ID,
		Key:       opts.key,
	})
	if err != nil {
		return fmt.Errorf("metadata key %q not found", opts.key)
	}

	entry := metadataEntry{
		Key:   m.Key,
		Value: m.Value,
	}

	formatted, err := utils.Format(entry, output)
	if err != nil {
		return err
	}

	fmt.Fprintln(out, formatted)
	return nil
}
