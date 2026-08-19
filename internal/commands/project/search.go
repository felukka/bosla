package project

import (
	"context"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/models"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type searchOptions struct {
	pattern string
}

func NewSearchCommand(q *database.Queries) *cobra.Command {
	opts := &searchOptions{}

	search := &cobra.Command{
		Use:   "search <regex>",
		Short: "Search projects with regex",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pattern = args[0]
		},
		RunE: utils.Wrap(func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
			projects, err := searchProjects(ctx, opts.pattern, q)
			if err != nil {
				return err
			}

			formatted, err := utils.Format(projects, outputFormat)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), formatted)
			return nil
		}),
	}

	return search
}

func searchProjects(ctx context.Context, pattern string, q *database.Queries) ([]models.Project, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern %q: %w", pattern, err)
	}

	allProjects, err := q.GetProjects(ctx)
	if err != nil {
		return nil, err
	}

	projects := models.ToProjects(allProjects)
	results := make([]models.Project, 0)

	for _, p := range projects {
		if re.MatchString(p.Name) ||
			re.MatchString(p.Status) ||
			re.MatchString(p.Link) ||
			re.MatchString(p.Description) {
			results = append(results, p)
		}
	}

	return results, nil
}
