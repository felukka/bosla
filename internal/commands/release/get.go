package release

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type getOptions struct {
	projectName string
	version     string
}

type releaseDetailsOutput struct {
	Project            string `json:"project"`
	ReleaseVersion     string `json:"release_version"`
	ReleaseDescription string `json:"release_description,omitempty"`
	ReleaseCreatedAt   string `json:"release_created_at"`
	Service            string `json:"service,omitempty"`
	ServiceVersion     string `json:"service_version,omitempty"`
	ServiceStatus      string `json:"service_status,omitempty"`
	ServiceGitHash     string `json:"service_git_hash,omitempty"`
}

func NewGetCommand(q *database.Queries) *cobra.Command {
	opts := &getOptions{}

	get := &cobra.Command{
		Use:   "get <project> <release-version>",
		Short: "Show one release with attached service versions",
		Args:  cobra.ExactArgs(2),
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.projectName = args[0]
			opts.version = args[1]
		},
		RunE: utils.Wrap(func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
			rows, err := getRelease(ctx, opts, q)
			if err != nil {
				return err
			}

			formatted, err := utils.Format(rows, outputFormat)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), formatted)
			return nil
		}),
	}

	return get
}

func getRelease(ctx context.Context, opts *getOptions, q *database.Queries) ([]releaseDetailsOutput, error) {
	projectID, err := q.GetProjectIdByName(ctx, opts.projectName)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", opts.projectName)
	}

	releaseRows, err := q.GetReleaseDetails(ctx, database.GetReleaseDetailsParams{
		ProjectID: projectID,
		Version:   opts.version,
	})
	if err != nil {
		return nil, err
	}
	if len(releaseRows) == 0 {
		return nil, fmt.Errorf(ErrReleaseNotFound, opts.version, opts.projectName)
	}

	results := make([]releaseDetailsOutput, 0, len(releaseRows))
	for _, row := range releaseRows {
		results = append(results, releaseDetailsOutput{
			Project:            row.ProjectName,
			ReleaseVersion:     row.ReleaseVersion,
			ReleaseDescription: row.ReleaseDescription,
			ReleaseCreatedAt:   row.ReleaseCreatedAt.Format(time.RFC3339),
			Service:            nullStringValue(row.ServiceName),
			ServiceVersion:     nullStringValue(row.ServiceVersion),
			ServiceStatus:      nullStringValue(row.ServiceStatus),
			ServiceGitHash:     nullStringValue(row.ServiceGitHash),
		})
	}

	return results, nil
}

func nullStringValue(v sql.NullString) string {
	if !v.Valid {
		return ""
	}
	return v.String
}
