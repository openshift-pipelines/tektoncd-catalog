package cmd

import (
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/catalog"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
	fc "github.com/openshift-pipelines/tektoncd-catalog/internal/fetcher/config"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/runner"
	"github.com/spf13/cobra"
)

// GenerateCmd represents the "generate" subcommand to generate the signature of a resource file.
type GenerateCmd struct {
	cmd    *cobra.Command // cobra command definition
	config string         // path for the catalog configuration file
	target string         // path to the folder where we want to generate the catalog
}

var _ runner.SubCommand = &GenerateCmd{}

const generateLongDescription = `# catalog-cd generate

Generates a file-based catalog in the target folder, based of a configuration file.

  $ catalog-cd generate \
      --config="/path/to/external.yaml" \
      /path/to/catalog/target
`

// Cmd exposes the cobra command instance.
func (v *GenerateCmd) Cmd() *cobra.Command {
	return v.cmd
}

// Complete asserts the required flags are informed, and the last argument is the resource file for
// signature verification.
func (v *GenerateCmd) Complete(_ *config.Config, args []string) error {
	if v.config == "" {
		return fmt.Errorf("flag --config is required")
	}

	if len(args) != 1 {
		return fmt.Errorf("you must specify a target to generate the catalog in")
	}
	v.target = args[0]
	return nil
}

// Validate asserts all the required files exists.
func (v *GenerateCmd) Validate() error {
	required := []string{
		v.config,
		// v.target,
	}
	for _, f := range required {
		if _, err := os.Stat(f); err != nil {
			return err
		}
	}
	return nil
}

// Run wrapper around "cosign generate-blob" command.
func (v *GenerateCmd) Run(cfg *config.Config) error {
	cfg.Infof("Generating a catalog from %s in %s\n", v.config, v.target)
	ghclient, err := api.DefaultRESTClient()
	if err != nil {
		return err
	}

	e, err := fc.LoadExternal(v.config)
	if err != nil {
		return err
	}
	c, err := catalog.FetchFromExternal(e, ghclient)
	if err != nil {
		return err
	}

	return catalog.GenerateFilesystem(v.target, c)
}

// NewGenerateCatalogCmd instantiates the "generate" subcommand.
func NewGenerateCatalogCmd() runner.SubCommand {
	v := &GenerateCmd{
		cmd: &cobra.Command{
			Use:          "generate-catalog",
			Args:         cobra.ExactArgs(1),
			Long:         generateLongDescription,
			Short:        "Verifies the resource file signature",
			SilenceUsage: true,
		},
	}

	f := v.cmd.PersistentFlags()
	f.StringVar(&v.config, "config", v.config, "path of the catalog configuration file")

	return v
}
