package cmd

import (
	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/runner"
	"github.com/spf13/cobra"
)

func CatalogCmd(cfg *config.Config) *cobra.Command {
	catalogCmd := &cobra.Command{
		Use:  "catalog",
		Long: `Catalog management commands.`,
	}

	catalogCmd.AddCommand(runner.NewRunner(cfg, NewCatalogGenerateCmd()).Cmd())
	catalogCmd.AddCommand(runner.NewRunner(cfg, NewCatalogGenerateFromExternalCmd()).Cmd())
	catalogCmd.AddCommand(runner.NewRunner(cfg, NewCatalogExternalsCmd()).Cmd())

	return catalogCmd
}
