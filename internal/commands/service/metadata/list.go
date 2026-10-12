package metadata

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type listOptions struct {
	pname string
	sname string
}

type metadataEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewListCommand() *cobra.Command {
	opts := &listOptions{}

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List service metadata",
		Args:    cobra.ExactArgs(2),
		Example: "list <project_name> <service_name>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
			opts.sname = args[1]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				return list(ctx, opts, queries, output, cmd.OutOrStdout())
			},
		),
	}

	return list
}

func list(ctx context.Context, opts *listOptions, q *database.Queries, output string, out io.Writer) error {
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

	metadata, err := q.GetServiceMetadata(ctx, svc.ID)
	if err != nil {
		return err
	}

	entries := make([]metadataEntry, 0, len(metadata))
	for _, m := range metadata {
		entries = append(entries, metadataEntry{
			Key:   m.Key,
			Value: m.Value,
		})
	}

	formatted, err := utils.Format(entries, output)
	if err != nil {
		return err
	}

	fmt.Fprintln(out, formatted)
	return nil
}
