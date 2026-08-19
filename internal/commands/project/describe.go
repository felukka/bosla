package project

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type statusOptions struct {
	name   string
	status string
}

func NewDescribeCommand(q *database.Queries) *cobra.Command {
	opts := &statusOptions{}
	var queries *database.Queries

	status := &cobra.Command{
		Use:     "status",
		Aliases: []string{"st"},
		Args:    cobra.RangeArgs(1, 2),
		Short:   "Update project status",
		PreRun: func(cmd *cobra.Command, args []string) {
			queries = database.Get()
			opts.name = args[0]

			if len(args) == 2 {
				opts.status = args[1]
			}
			opts.name = args[0]
		},
	}

	status.RunE = func(cmd *cobra.Command, args []string) error {
		outputFormat, err := cmd.Flags().GetString("output")
		ctx := cmd.Context()

		if len(args) == 1 {
			s, err := getProjectStatus(ctx, opts, queries)
			if err != nil {
				return err
			}

			var fmtS string
			switch {
			case jsonFlag:
				fmtS, err = utils.FormatJSON(s)
			case yamlFlag:
				fmtS, err = utils.FormatTable(s)
			default:
				fmtS, err = utils.Format(s)
			}
			if err != nil {
				return err
			}

			fmt.Println(fmtS)

			return nil
		}

		if err := updateProjectStatus(ctx, opts, queries); err != nil {
			return err
		}

		return nil
	}

	return status
}

func updateProjectStatus(ctx context.Context, opts *statusOptions, q *database.Queries) error {
	if _, err := q.CheckProjectExistsByName(ctx, opts.name); err != nil {
		return fmt.Errorf(ErrProjectNotFound, opts.name)
	}

	if err := q.UpdateProjectStatusByName(ctx, database.UpdateProjectStatusByNameParams{
		Name:      opts.name,
		Status:    opts.status,
		UpdatedAt: time.Now().UTC(),
	}); err != nil {
		return err
	}

	return nil
}

func getProjectStatus(ctx context.Context, opts *statusOptions, q *database.Queries) (any, error) {
	pId, err := q.GetProjectIdByName(ctx, opts.name)
	if err != nil {
		return "", fmt.Errorf(ErrProjectNotFound, opts.name)
	}

	opts.status, err = q.GetProjectStatusById(ctx, pId)
	if err != nil {
		return "", fmt.Errorf(ErrProjectNotFound, err)
	}

	return struct {
		Project string `json:"project"`
		Status  string `json:"status"`
	}{
		Project: opts.name,
		Status:  opts.status,
	}, nil
}
