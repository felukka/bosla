package release

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/commands/project"
	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type listOptions struct {
	pname string
}

type releaseListOutput struct {
	Project     string `json:"project"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
}

func NewListCommand(q *database.Queries) *cobra.Command {
	opts := &listOptions{}

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List releases for a project",
		Args:    cobra.ExactArgs(1),
		Example: "list <project_name>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
				releases, err := listReleases(ctx, opts.pname, q)
				if err != nil {
					return err
				}

				formatted, err := utils.Format(releases, outputFormat)
				if err != nil {
					return err
				}

				if _, err := fmt.Fprintln(cmd.OutOrStdout(), formatted); err != nil {
					return err
				}

				return nil
			},
		),
	}

	return list
}

func listReleases(
	ctx context.Context,
	projectName string,
	q *database.Queries,
) ([]releaseListOutput, error) {
	projectID, err := q.GetProjectIdByName(ctx, projectName)
	if err != nil {
		return nil, fmt.Errorf(project.ErrProjectNotFound, projectName)
	}

	projectVersions, err := q.ListProjectVersions(ctx, projectID)
	if err != nil {
		return nil, err
	}

	results := make([]releaseListOutput, 0, len(projectVersions))
	for _, item := range projectVersions {
		results = append(results, releaseListOutput{
			Project:     projectName,
			Version:     item.Version,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		})
	}

	return results, nil
}
