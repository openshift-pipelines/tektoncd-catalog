package runner

import (
	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/config"

	"github.com/spf13/cobra"
)

// Runner controls the SubCommand workflow from end-to-end.
type Runner struct {
	cfg        *config.Config // global configuration
	subCommand SubCommand     // SubCommand instance
}

// Cmd exposes the subcommand's cobra command instance.
func (r *Runner) Cmd() *cobra.Command {
	return r.subCommand.Cmd()
}

// RunE cobra's RunE function, executes the whole SubCommand workflow.
func (r *Runner) RunE(_ *cobra.Command, args []string) error {
	err := r.subCommand.Complete(r.cfg, args)
	if err != nil {
		return err
	}
	if err = r.subCommand.Validate(); err != nil {
		return err
	}
	return r.subCommand.Run(r.cfg)
}

// NewRunner instantiates a Runner, making sure it's RunE command is mapped to the local method,
// which executes the whole interface workflow.
func NewRunner(cfg *config.Config, subCommand SubCommand) *Runner {
	r := &Runner{
		cfg:        cfg,
		subCommand: subCommand,
	}
	r.subCommand.Cmd().RunE = r.RunE
	return r
}
