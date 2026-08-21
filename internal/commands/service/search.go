package service

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
	pname   string
	pattern string
}

func NewSearchCommand() *cobra.Command {
	opts := &searchOptions{}

	search := &cobra.Command{
		Use:     "search",
		Short:   "Search services in a project with a pattern",
		Args:    cobra.ExactArgs(2),
		Example: "search <project_name> <pattern>",
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
			opts.pattern = args[1]
		},
		RunE: utils.Wrap(
			func(ctx context.Context, cmd *cobra.Command, args []string, queries *database.Queries, output string) error {
				services, err := search(ctx, opts, queries)
				if err != nil {
					return err
				}

				output, err = utils.Format(services, output)
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
	opts *searchOptions,
	q *database.Queries,
) ([]models.Service, error) {
	re, err := regexp.Compile(opts.pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern %q: %w", opts.pattern, err)
	}

	projectID, err := q.GetProjectIdByName(ctx, opts.pname)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", opts.pname)
	}

	allServices, err := q.GetServiceByProjectName(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf(ErrServiceList, opts.pname, err)
	}

	services := models.ToServices(allServices)
	results := make([]models.Service, 0)

	for _, s := range services {
		if re.MatchString(s.Name) ||
			re.MatchString(s.Status) ||
			re.MatchString(s.RepoURL) ||
			re.MatchString(s.Description) {
			results = append(results, s)
		}
	}

	return results, nil
}
