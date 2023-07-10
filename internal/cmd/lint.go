package cmd

import (
	"fmt"
	"os"

	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/linter"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/runner"

	"github.com/spf13/cobra"
)

// LintCmd "lint" subcommand, to assert the best practices on a Tekton resource file.
type LintCmd struct {
	cmd      *cobra.Command // cobra command definition
	resource string         // path to the resource file
}

var _ runner.SubCommand = &LintCmd{}

const lintLongDescription = `# catalog-cd lint

Lints the informed resource file informing which attributes aren't following the best practices.
`

// Cmd shares the Cobra command instance.
func (l *LintCmd) Cmd() *cobra.Command {
	return l.cmd
}

// Complete asserts a single argument is informed.
func (l *LintCmd) Complete(_ *config.Config, args []string) error {
	if len(args) == 1 {
		l.resource = args[0]
		return nil
	}
	return fmt.Errorf("you must inform a single argument (%d)", len(args))
}

// Validate assert the informed file exists.
func (l *LintCmd) Validate() error {
	_, err := os.Stat(l.resource)
	return err
}

// Run runs the linter against the informed resource file.
func (l *LintCmd) Run(cfg *config.Config) error {
	lr, err := linter.NewLinter(cfg, l.resource)
	if err != nil {
		return err
	}
	return lr.Enforce()
}

// NewLintCmd instantiates the "lint" subcommand.
func NewLintCmd() runner.SubCommand {
	cmd := &cobra.Command{
		Use:          "lint",
		Short:        "Tekton resource file linter",
		Long:         lintLongDescription,
		SilenceUsage: true,
	}
	return &LintCmd{cmd: cmd}
}
