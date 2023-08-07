package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

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
	isFolder, err := isDirectory(l.resource)
	if err != nil {
		return err
	}
	if isFolder {
		var wg sync.WaitGroup
		var mu sync.Mutex
		errs := []error{}

		err := filepath.Walk(l.resource,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !isYaml(path) {
					fmt.Println("Ignoring non-yaml file", path)
					return nil
				}
				wg.Add(1)
				go func() {
					err := lintFile(cfg, path)
					mu.Lock()
					errs = append(errs, err)
					mu.Unlock()
					wg.Done()
				}()
				return nil
			})
		if err != nil {
			return err
		}
		wg.Wait()
		return errors.Join(errs...)
	}
	return lintFile(cfg, l.resource)
}

func lintFile(cfg *config.Config, path string) error {
	lr, err := linter.NewLinter(cfg, path)
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

// isDirectory determines if a file represented
// by `path` is a directory or not
func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func isYaml(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".yaml" || ext == ".yml"
}
