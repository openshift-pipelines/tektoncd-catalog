package cmd

import (
	"archive/tar"
	"compress/gzip"
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
	fmt.Fprintf(os.Stderr, "# Found %d path to inspect!\n", len(r.paths))
	return nil
}

// Run creates a ".catalog.yaml" (contract file) with the release scope, saves the contract
// on the location informed by the "--output" flag.
func (r *ReleaseCmd) Run(_ *config.Config) error {
	c := contract.NewContractEmpty()
	// going through the pattern slice collected before to select the tekton resource files
	// to be part of the current release, in other words, release scope
	fmt.Fprintf(os.Stderr, "# Scan Tekton resources on: %s\n", strings.Join(r.paths, ", "))
	for _, p := range r.paths {
		files, err := resource.Scanner(p)
		if err != nil {
			return err
		}

		for _, f := range files {
			fmt.Fprintf(os.Stderr, "# Loading resource file: %q\n", f)
			taskname := filepath.Base(filepath.Dir(f))
			resourceType, err := resource.GetResourceType(f)
			if err != nil {
				return err
			}
			resourceFolder := filepath.Join(r.output, strings.ToLower(resourceType)+"s", taskname)
			if err := os.MkdirAll(resourceFolder, os.ModePerm); err != nil {
				return err
			}
			if err := c.AddResourceFile(f, r.version); err != nil {
				if errors.Is(err, contract.ErrTektonResourceUnsupported) {
					return err
				}
				fmt.Printf("# WARNING: Skipping file %q!\n", f)
			}
			// Copy it to output
			if err := copyFile(f, filepath.Join(resourceFolder, filepath.Base(f))); err != nil {
				return err
			}
			readmeFile := filepath.Join(filepath.Dir(f), "README.md")
			if _, err := os.Stat(readmeFile); err == nil {
				// This is the README, copy it to output
				if err := copyFile(readmeFile, filepath.Join(resourceFolder, "README.md")); err != nil {
					return err
				}
				continue
			}
		}
	}

	catalogPath := filepath.Join(r.output, "catalog.yaml")
	fmt.Fprintf(os.Stderr, "# Saving release contract at %q\n", catalogPath)
	if err := c.SaveAs(catalogPath); err != nil {
		return err
	}

	// Create a tarball (without catalog.yaml
	tarball := filepath.Join(r.output, "resources.tar.gz")
	fmt.Fprintf(os.Stderr, "# Creating tarball at %q\n", tarball)
	if err := createTektonResourceArchive(tarball, r.output); err != nil {
		return err
	}
	return nil
}

func createTektonResourceArchive(archiveFile, output string) error {
	// Create output file
	out, err := os.Create(archiveFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Create the archive
	return createArchive(output, out)
}
func createArchive(output string, buf io.Writer) error {
	// Create new Writers for gzip and tar
	// These writers are chained. Writing to the tar writer will
	// write to the gzip writer which in turn will write to
	// the "buf" writer
	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Iterate over files and add them to the tar archive
	return filepath.Walk(output, func(file string, fi os.FileInfo, err error) error {
		// return on any error
		if err != nil {
			return err
		}
		if filepath.Base(file) == "catalog.yaml" || filepath.Base(file) == "resources.tar.gz" {
			return nil
		}
		if fi.IsDir() || !fi.Mode().IsRegular() {
			return nil
		}
		return addToArchive(tw, file, output)
	})
}

func addToArchive(tw *tar.Writer, filename string, output string) error {
	// Open the file which will be written into the archive
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get FileInfo about our file providing file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory strucuture would
	// not be preserved
	// https://golang.org/src/archive/tar/common.go?#L626
	header.Name = strings.TrimPrefix(filename, filepath.Base(output)+"/")

	// Write file header to the tar archive
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// Copy file content to tar archive
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
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
