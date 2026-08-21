package configure

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	ConfigFile = func() string {
		home, err := os.UserHomeDir()
		if err != nil {
			return ".bosla.json"
		}

		return filepath.Join(home, ".bosla.json")
	}()
	ErrConfigInValid = errors.New("bosla.json is invalid")
)

func NewConfigureCommand() *cobra.Command {
	configure := &cobra.Command{
		Use:     "configure",
		Aliases: []string{"config"},
		Short:   "manage configuration for bosla",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return nil
			}

			return nil
		},
	}

	configure.AddCommand(
		NewInitCommand(),
		NewSetCommand(),
	)

	return configure
}
