package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/models"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type searchOptions struct {
	query string
}

func NewSearchCommand(q *database.Queries) *cobra.Command {
	opts := &searchOptions{}

	search := &cobra.Command{
		Use:   "search <query>",
		Short: "Search projects",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.query = strings.ToLower(args[0])
		},
		RunE: utils.Wrap(func(ctx context.Context, cmd *cobra.Command, args []string, outputFormat string) error {
			projects, err := searchProjects(ctx, opts.query, q)
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

func searchProjects(ctx context.Context, query string, q *database.Queries) ([]models.Project, error) {
	allProjects, err := q.GetProjects(ctx)
	if err != nil {
		return nil, err
	}

	projects := models.ToProjects(allProjects)
	results := make([]models.Project, 0)

	for _, p := range projects {
		if strings.Contains(strings.ToLower(p.Name), query) ||
			strings.Contains(strings.ToLower(p.Status), query) ||
			strings.Contains(strings.ToLower(p.Link), query) ||
			strings.Contains(strings.ToLower(p.Description), query) {
			results = append(results, p)
		}
	}

	return results, nil
}
