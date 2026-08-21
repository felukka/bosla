package configure

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

type setOptions struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	URL      string `json:"url"`
	Port     uint16 `json:"port"`
	SkipSSL  bool   `json:"skipSsl"`
}

func NewSetCommand() *cobra.Command {
	opts := &setOptions{}

	set := &cobra.Command{
		Use:   "set",
		Short: "Manage configuration for bosla",
		Example: `
bosla configure set -u <username> -p <password> [--url <url|127.0.0.1>]
OR
export BOSLA_DATABASE_URL="postgres://user:password@localhost:5432/bosla"
AND
bosla configure init
`,
		PersistentPreRunE:  func(cmd *cobra.Command, args []string) error { return nil },
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error { return nil },
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := set(opts); err != nil {
				return err
			}

			return nil
		},
	}

	set.Flags().StringVarP(&opts.Username, "username", "u", "", "username of database")
	set.Flags().StringVarP(&opts.Password, "password", "p", "", "password of user")
	set.Flags().StringVarP(&opts.Database, "database", "d", "bosla", "password of user")
	set.Flags().StringVar(&opts.URL, "url", "127.0.0.1", "url of database")
	set.Flags().Uint16Var(&opts.Port, "port", 5432, "port of database")
	set.Flags().
		BoolVar(&opts.SkipSSL, "skip-ssl", false, "disable ssl when connecting to database")

	return set
}

func set(opts *setOptions) error {
	ecfg, err := LoadConfig()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, ErrConfigInValid) {
			save(opts)
			return nil
		}

		return err
	}

	if opts.Username != "" {
		ecfg.Username = opts.Username
	}
	if opts.Password != "" {
		ecfg.Password = opts.Password
	}
	if opts.URL != "" && opts.URL != "127.0.0.1" {
		ecfg.URL = opts.URL
	}

	if opts.Port != 0 {
		ecfg.Port = opts.Port
	}

	ecfg.SkipSSL = opts.SkipSSL

	save(ecfg)

	return nil
}

func save(cfg *setOptions) {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	_ = os.WriteFile(ConfigFile, data, 0644)
}

func IsValidConfig() error {
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		return err
	}

	var cfg setOptions
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	if cfg.URL == "" || cfg.Username == "" {
		return ErrConfigInValid
	}

	return nil
}

func LoadConfig() (*setOptions, error) {
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		return nil, err
	}

	var cfg setOptions
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.URL == "" || cfg.Username == "" {
		return nil, ErrConfigInValid
	}

	return &cfg, nil
}
