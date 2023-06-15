package cmd

import (
	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/config"
	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/runner"

	"github.com/spf13/cobra"
	tkncli "github.com/tektoncd/cli/pkg/cli"
)

func NewRootCmd(stream *tkncli.Stream) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "catalog-cd",
		Long: `TODO`,
	}

	cfg := config.NewConfigWithFlags(stream, rootCmd.PersistentFlags())

	rootCmd.AddCommand(runner.NewRunner(cfg, NewLintCmd()).Cmd())
	rootCmd.AddCommand(runner.NewRunner(cfg, NewProbeCmd()).Cmd())
	rootCmd.AddCommand(runner.NewRunner(cfg, NewRenderCmd()).Cmd())
	rootCmd.AddCommand(runner.NewRunner(cfg, NewVerifyCmd()).Cmd())

	return rootCmd
}
