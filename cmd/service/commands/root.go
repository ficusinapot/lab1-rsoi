package commands

import (
	"os"

	serviceconfig "github.com/ficusinapot/ds/cmd/service/config"
	"github.com/ficusinapot/ds/cmd/service/errors"
	"github.com/ficusinapot/ds/cmd/service/version"

	"github.com/spf13/cobra"
)

type RunFunc func(cfgPath string) error

func Execute(run RunFunc) error {
	rootCmd := &cobra.Command{ //nolint:exhaustruct_v5
		Use:          "service",
		Short:        "Run service",
		Version:      version.GetApplicationVersion(),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = args

			return run(configPath(cmd))
		},
	}

	rootCmd.PersistentFlags().StringP("config", "c", serviceconfig.DefaultConfigPath, "config file")

	if err := rootCmd.Execute(); err != nil {
		return errors.CommandLineProcessingError.Wrap(err, "CLI error")
	}

	return nil
}

func configPath(cmd *cobra.Command) string {
	flag := cmd.Flag("config")
	if flag != nil && flag.Changed {
		return flag.Value.String()
	}

	if envPath := os.Getenv(serviceconfig.EnvConfigPathKey); envPath != "" {
		return envPath
	}

	return serviceconfig.DefaultConfigPath
}
