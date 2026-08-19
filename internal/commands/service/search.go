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
	projectName string
	pattern     string
}

func NewSearchCommand(q *database.Queries) *cobra.Command {
	opts := &searchOptions{}

	search := &cobra.Command{
		Use:   "search <project> <regex>",
		Short: "Search services in a project with regex",
		Args:  cobra.ExactArgs(2),
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.projectName = args[0]
			opts.pattern = args[1]
		},
		RunE: utils.Wrap(func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
			services, err := searchServices(ctx, opts, q)
			if err != nil {
				return err
			}

			formatted, err := utils.Format(services, outputFormat)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), formatted)
			return nil
		}),
	}

	return search
}

func searchServices(ctx context.Context, opts *searchOptions, q *database.Queries) ([]models.Service, error) {
	re, err := regexp.Compile(opts.pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern %q: %w", opts.pattern, err)
	}

	projectID, err := q.GetProjectIdByName(ctx, opts.projectName)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", opts.projectName)
	}

	allServices, err := q.GetServiceByProjectName(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf(ErrServiceList, opts.projectName, err)
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
