package project

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
)

type metadataOptions struct {
	projectName string
	key         string
	set         bool
	get         bool
	delete      bool
}

func NewMetadataCommand() *cobra.Command {
	opts := &metadataOptions{}
	var queries *database.Queries

	metadata := &cobra.Command{
		Use:     "metadata [project-name]",
		Aliases: []string{"md"},
		Short:   "Manage project metadata",
		Args:    cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = database.Get()
			opts.projectName = args[0]
		},
	}

	flags := metadata.Flags()
	flags.StringP("get", "g", "", "Get metadata value for the specified key")
	flags.StringP(
		"set",
		"s",
		"",
		"Set metadata value for the specified key in the format key=value",
	)
	flags.StringP("delete", "d", "", "Delete metadata entry for the specified key")

	metadata.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		ctx := cmd.Context()

		get, _ := flags.GetString("get")
		set, _ := flags.GetString("set")
		del, _ := flags.GetString("delete")

		if get != "" {
			opts.key = get
			mt, err := getProjectMetadata(ctx, opts, queries)
			if err != nil {
				return fmt.Errorf("failed to get project metadata: %w", err)
			}
			fmt.Println(mt)
			return nil
		}
		if set != "" {
			fmt.Println("set...")
			return nil
		}
		if del != "" {
			fmt.Println("delete...")
			return nil
		}

		mt, err := getProjectMetadata(ctx, opts, queries)
		if err != nil {
			return fmt.Errorf("failed to get project metadata: %w", err)
		}
		fmt.Println(mt)
		return nil
	}

	return metadata
}

func getProjectMetadata(
	ctx context.Context,
	opts *metadataOptions,
	q *database.Queries,
) (map[string]any, error) {
	exist, err := q.CheckProjectExistsByName(ctx, opts.projectName)
	if err != nil {
		return nil, fmt.Errorf("failed to check project existence: %w", err)
	}
	if !exist {
		return nil, fmt.Errorf("project %s does not exist", opts.projectName)
	}

	rawMd, err := q.GetProjectMetadata(ctx, opts.projectName)
	if err != nil {
		return nil, fmt.Errorf("project has no metadata: %w", err)
	}

	var md map[string]any
	if err := json.Unmarshal(rawMd, &md); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project metadata: %w", err)
	}

	if opts.key != "" {
		value, ok := md[opts.key]
		if !ok {
			return nil, fmt.Errorf(keyNotFoundInMetadataErr, opts.key)
		}
		return map[string]any{opts.key: value}, nil
	}

	return md, nil
}
