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

func NewSearchCommand() *cobra.Command {
	opts := &searchOptions{}

	search := &cobra.Command{
		Use:     "search",
		Short:   "Search projects with a pattern",
		Args:    cobra.ExactArgs(1),
		Example: "search <pattern>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pattern = args[0]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				projects, err := search(ctx, opts.pattern, queries)
				if err != nil {
					return err
				}

				output, err = utils.Format(projects, output)
				if err != nil {
					return err
				}

				if _, err := fmt.Fprintln(cmd.OutOrStdout(), output); err != nil {
					return err
				}
				return nil
			},
		),
	}

	return search
}

func search(
	ctx context.Context,
	pattern string,
	q *database.Queries,
) ([]models.Project, error) {
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
		if re.MatchString(p.Name) {
			results = append(results, p)
		}
	}

	return results, nil
}
