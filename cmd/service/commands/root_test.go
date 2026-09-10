package commands

import (
	"testing"

	serviceconfig "github.com/ficusinapot/ds/cmd/service/config"

	"github.com/spf13/cobra"
)

func TestConfigPathUsesDefaultPath(t *testing.T) {
	t.Setenv(serviceconfig.EnvConfigPathKey, "")

	cmd := commandWithConfigFlag(t, serviceconfig.DefaultConfigPath)

	if got := configPath(cmd); got != serviceconfig.DefaultConfigPath {
		t.Fatalf("unexpected config path: %q", got)
	}
}

func TestConfigPathUsesEnvironmentPath(t *testing.T) {
	const envPath = "/tmp/service.yaml"

	t.Setenv(serviceconfig.EnvConfigPathKey, envPath)

	cmd := commandWithConfigFlag(t, serviceconfig.DefaultConfigPath)

	if got := configPath(cmd); got != envPath {
		t.Fatalf("unexpected config path: %q", got)
	}
}

func TestConfigPathPrefersFlagOverEnvironment(t *testing.T) {
	const flagPath = "/tmp/flag.yaml"

	t.Setenv(serviceconfig.EnvConfigPathKey, "/tmp/env.yaml")

	cmd := commandWithConfigFlag(t, serviceconfig.DefaultConfigPath)
	if err := cmd.Flags().Set("config", flagPath); err != nil {
		t.Fatalf("set config flag: %v", err)
	}

	if got := configPath(cmd); got != flagPath {
		t.Fatalf("unexpected config path: %q", got)
	}
}

func commandWithConfigFlag(t *testing.T, defaultValue string) *cobra.Command {
	t.Helper()

	cmd := &cobra.Command{} //nolint:exhaustruct_v5
	cmd.Flags().StringP("config", "c", defaultValue, "config file")

	return cmd
}
