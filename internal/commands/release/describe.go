package release

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type describeOptions struct {
	pname   string
	version string
}

type releaseDescribeOutput struct {
	Project        string `json:"project"`
	Version        string `json:"version"`
	Description    string `json:"description,omitempty"`
	CreatedAt      string `json:"created_at"`
	ServiceName    string `json:"service_name,omitempty"`
	ServiceVersion string `json:"service_version,omitempty"`
	ServiceStatus  string `json:"service_status,omitempty"`
	ServiceGitHash string `json:"service_git_hash,omitempty"`
}

func NewDescribeCommand() *cobra.Command {
	opts := &describeOptions{}

	describe := &cobra.Command{
		Use:     "describe",
		Args:    cobra.ExactArgs(2),
		Short:   "Describe a release",
		Example: "describe <project_name> <version>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
			opts.version = args[1]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				releases, err := describe(ctx, opts, queries)
				if err != nil {
					return err
				}

				formatted, err := utils.Format(releases, output)
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

	return describe
}

func describe(
	ctx context.Context,
	opts *describeOptions,
	q *database.Queries,
) ([]releaseDescribeOutput, error) {
	projectID, err := q.GetProjectIdByName(ctx, opts.pname)
	if err != nil {
		return nil, fmt.Errorf(project.ErrProjectNotFound, opts.pname)
	}

	rows, err := q.GetReleaseDetails(ctx, database.GetReleaseDetailsParams{
		ProjectID: projectID,
		Version:   opts.version,
	})
	if err != nil {
		return nil, err
	}

	results := make([]releaseDescribeOutput, 0, len(rows))
	for _, row := range rows {
		item := releaseDescribeOutput{
			Project:     row.ProjectName,
			Version:     row.ReleaseVersion,
			Description: row.ReleaseDescription,
			CreatedAt:   row.ReleaseCreatedAt.Format("2006-01-02T15:04:05Z"),
		}

		if row.ServiceName.Valid {
			item.ServiceName = row.ServiceName.String
		}
		if row.ServiceVersion.Valid {
			item.ServiceVersion = row.ServiceVersion.String
		}
		if row.ServiceStatus.Valid {
			item.ServiceStatus = row.ServiceStatus.String
		}
		if row.ServiceGitHash.Valid {
			item.ServiceGitHash = row.ServiceGitHash.String
		}

		results = append(results, item)
	}

	return results, nil
}
