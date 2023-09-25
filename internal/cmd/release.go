package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/contract"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/resource"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/runner"

	"github.com/spf13/cobra"
)

// ReleaseCmd creates a contract (".catalog.yaml") based on Tekton resources files.
type ReleaseCmd struct {
	cmd     *cobra.Command // cobra command definition
	version string         // release version
	files   []string       // tekton resource files
	output  string         // output file, where the contract will be written
}

var _ runner.SubCommand = &ReleaseCmd{}

const releaseLongDescription = `# catalog-cd release

Creates a contract file (".catalog.yaml") for the Tekton resources specified on
the last argument(s), the contract is stored on the "--output" location, or by
default ".catalog.yaml" on the current directory.

The following examples will store the ".catalog.yaml" on the current directory, in
order to change its location see "--output" flag.

  # release all "*.yaml" files on the subdirectory
  $ catalog-cd release --version="0.0.1" path/to/tekton/files/*.yaml

  # release all "*.{yml|yaml}" files on the current directory
  $ catalog-cd release --version="0.0.1" *.yml *.yaml

  # release all "*.yml" and "*.yaml" files from the current directory
  $ catalog-cd release --version="0.0.1"

It always require the "--version" flag specifying the common revision for all
resources in scope.
`

// Cmd exposes the cobra command instance.
func (r *ReleaseCmd) Cmd() *cobra.Command {
	return r.cmd
}

// gatherGlogPatternFromArgs creates a set of glob patterns based on the final command line
// arguments (args), when the slice is empty it assumes the current working directory.
func (*ReleaseCmd) gatherGlogPatternFromArgs(args []string) ([]string, error) {
	patterns := []string{}

	if len(args) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, wd)
		fmt.Printf("# Using current directory: %q\n", wd)
	} else {
		patterns = append(patterns, args...)
	}

	return patterns, nil
}

// Complete creates the "release" scope by finding all Tekton resource files using the cli
// args glob pattern(s).
func (r *ReleaseCmd) Complete(_ *config.Config, args []string) error {
	// making sure the output flag is informed before attempt to search files
	if r.output == "" {
		return fmt.Errorf("--output flag is not informed")
	}

	// putting together a slice of glob patterns to search tekton files
	patterns, err := r.gatherGlogPatternFromArgs(args)
	if err != nil {
		return err
	}

	// going through the pattern slice collected before to select the tekton resource files
	// to be part of the current release, in other words, release scope
	fmt.Printf("# Scan Tekton resources on: %s\n", strings.Join(patterns, ", "))
	for _, pattern := range patterns {
		files, err := resource.Scanner(pattern)
		if err != nil {
			return err
		}
		r.files = append(r.files, files...)
	}
	return nil
}

// Validate assert the release scope is not empty.
func (r *ReleaseCmd) Validate() error {
	if len(r.files) == 0 {
		return fmt.Errorf("no tekton resource files have been found")
	}
	fmt.Printf("# Found %d files to inspect!\n", len(r.files))
	return nil
}

// Run creates a ".catalog.yaml" (contract file) with the release scope, saves the contract
// on the location informed by the "--output" flag.
func (r *ReleaseCmd) Run(_ *config.Config) error {
	c := contract.NewContractEmpty()

	fmt.Printf("# Generating contract for release %q...\n", r.version)
	for _, f := range r.files {
		fmt.Printf("# Loading resource file: %q\n", f)
		if err := c.AddResourceFile(f, r.version); err != nil {
			if errors.Is(err, contract.ErrTektonResourceUnsupported) {
				return err
			}
			fmt.Printf("# WARNING: Skipping file %q!\n", f)
		}
	}

	fmt.Printf("# Saving release contract at %q\n", r.output)
	return c.SaveAs(r.output)
}

// NewReleaseCmd instantiates the NewReleaseCmd subcommand and flags.
func NewReleaseCmd() runner.SubCommand {
	r := &ReleaseCmd{
		cmd: &cobra.Command{
			Use:          "release [flags] [glob|directory]",
			Short:        "Creates a contract for Tekton resource files",
			Long:         releaseLongDescription,
			Args:         cobra.ArbitraryArgs,
			SilenceUsage: true,
		},
	}

	f := r.cmd.PersistentFlags()

	f.StringVar(&r.version, "version", "", "release version")
	f.StringVar(&r.output, "output", contract.Filename, "path to the contract file")

	if err := r.cmd.MarkPersistentFlagRequired("version"); err != nil {
		panic(err)
	}

	return r
}
