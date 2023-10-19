package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	paths   []string       // tekton resource paths
	output  string         // output path, where the contract and tarball will be written
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

// Complete creates the "release" scope by finding all Tekton resource files using the cli
// args glob pattern(s).
func (r *ReleaseCmd) Complete(_ *config.Config, args []string) error {
	// making sure the output flag is informed before attempt to search files
	if r.output == "" {
		return fmt.Errorf("--output flag is not informed")
	}
	r.paths = args
	return nil
}

// Validate assert the release scope is not empty.
func (r *ReleaseCmd) Validate() error {
	if len(r.paths) == 0 {
		return fmt.Errorf("no tekton resource paths have been found")
	}
	fmt.Printf("# Found %d path to inspect!\n", len(r.paths))
	return nil
}

// Run creates a ".catalog.yaml" (contract file) with the release scope, saves the contract
// on the location informed by the "--output" flag.
func (r *ReleaseCmd) Run(_ *config.Config) error {
	c := contract.NewContractEmpty()
	// going through the pattern slice collected before to select the tekton resource files
	// to be part of the current release, in other words, release scope
	fmt.Printf("# Scan Tekton resources on: %s\n", strings.Join(r.paths, ", "))
	for _, p := range r.paths {
		files, err := resource.Scanner(p)
		if err != nil {
			return err
		}

		for _, f := range files {
			fmt.Fprintf(os.Stderr, "# Loading resource file: %q\n", f)
			taskname := filepath.Base(filepath.Dir(f))
			if filepath.Base(f) == "README.md" {
				// This is the README, copy it to output
				if err := os.MkdirAll(filepath.Join(r.output, taskname), os.ModePerm); err != nil {
					return err
				}
				if err := copyFile(f, filepath.Join(r.output, taskname, "README.md")); err != nil {
					return err
				}
				continue
			}
			if err := c.AddResourceFile(f, r.version); err != nil {
				if errors.Is(err, contract.ErrTektonResourceUnsupported) {
					return err
				}
				fmt.Printf("# WARNING: Skipping file %q!\n", f)
			}

			// Copy it to output
			if err := copyFile(f, filepath.Join(r.output, taskname, filepath.Base(f))); err != nil {
				return err
			}
		}
	}

	catalogPath := filepath.Join(r.output, "catalog.yaml")
	fmt.Printf("# Saving release contract at %q\n", catalogPath)
	return c.SaveAs(catalogPath)
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
	f.StringVar(&r.output, "output", contract.Filename, "path to the release files (to attach to a given release)")

	if err := r.cmd.MarkPersistentFlagRequired("version"); err != nil {
		panic(err)
	}

	return r
}

func copyFile(src, dst string) error {
	// Open the source file for reading
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Flush the destination file to ensure all data is written
	err = dstFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
