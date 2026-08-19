package service

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mahmoudk1000/bosla/internal/database"
	"github.com/mahmoudk1000/bosla/internal/models"
	"github.com/mahmoudk1000/bosla/internal/utils"
)

type getOptions struct {
	pname string
}

func NewListCommand(q *database.Queries) *cobra.Command {
	opts := &getOptions{}

	return &cobra.Command{
		Use:     "get <project>",
		Aliases: []string{"ls"},
		Short:   "Get/List all service of a project",
		Args:    cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			opts.pname = args[0]
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFormat, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}

			service, err := getService(opts, q)
			if err != nil {
				return err
			}

			output, err := utils.Format(service, outputFormat)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), output)

			return nil
		},
	}
}

func getService(
	opts getOptions,
	q *database.Queries,
) ([]models.Service, error) {
	ps, err := q.GetServiceByProjectName(ctx, opts.pname)
	if err != nil {
		return nil, fmt.Errorf(ErrServiceList, opts.pname, err)
	}

	return models.ToServices(ps), nil
}
