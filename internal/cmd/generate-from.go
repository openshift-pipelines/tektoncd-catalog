package cmd

import (
	"fmt"
	"path"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/catalog"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
	fc "github.com/openshift-pipelines/tektoncd-catalog/internal/fetcher/config"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/runner"
	"github.com/spf13/cobra"
)

// GenerateFromExternalCmd represents the "generate" subcommand to generate the signature of a resource file.
type GenerateFromExternalCmd struct {
	cmd            *cobra.Command // cobra command definition
	name           string         // name of the repository to pull (a bit useless)
	url            string         // url of the repository to pull
	resourceType   string         // type of resource to pull
	ignoreVersions string         // versions to ignore while pulling
	target         string         // path to the folder where we want to generate the catalog
}

var _ runner.SubCommand = &GenerateFromExternalCmd{}

const generateLongFromExternalDescription = `# catalog-cd generate

Generates a file-based catalog in the target folder, based of a configuration file.

  $ catalog-cd generate \
      --config="/path/to/external.yaml" \
      /path/to/catalog/target
`

// Cmd exposes the cobra command instance.
func (v *GenerateFromExternalCmd) Cmd() *cobra.Command {
	return v.cmd
}

// Complete asserts the required flags are informed, and the last argument is the resource file for
// signature verification.
func (v *GenerateFromExternalCmd) Complete(_ *config.Config, args []string) error {
	if v.url == "" {
		return fmt.Errorf("flag --config is required")
	}
	if v.resourceType == "" {
		return fmt.Errorf("flag --resourceType is required")
	}

	if len(args) != 1 {
		return fmt.Errorf("you must specify a target to generate the catalog in.")
	}
	v.target = args[0]
	return nil
}

// Validate asserts all the required files exists.
func (v *GenerateFromExternalCmd) Validate() error {
	return nil
}

// Run wrapper around "cosign generate-blob" command.
func (v *GenerateFromExternalCmd) Run(cfg *config.Config) error {
	cfg.Infof("Generating a partial catalog from %s (type: %s)\n", v.url, v.resourceType)
	ghclient, err := api.DefaultRESTClient()
	if err != nil {
		return err
	}

	name := v.name
	if name == "" {
		name = path.Base(v.url)
	}
	ignoreVersions := []string{}
	if v.ignoreVersions != "" {
		ignoreVersions = strings.Split(v.ignoreVersions, ",")
	}

	e := fc.External{
		Repositories: []fc.Repository{{
			Name:           name,
			URL:            v.url,
			Types:          []string{v.resourceType},
			IgnoreVersions: ignoreVersions,
		}},
	}
	c, err := catalog.FetchFromExternals(e, ghclient)
	if err != nil {
		return err
	}

	return catalog.GenerateFilesystem(v.target, c)
}

// NewCatalogGenerateFromExternalCmd instantiates the "generate" subcommand.
func NewCatalogGenerateFromExternalCmd() runner.SubCommand {
	v := &GenerateFromExternalCmd{
		cmd: &cobra.Command{
			Use:          "generate-from",
			Args:         cobra.ExactArgs(1),
			Long:         generateLongFromExternalDescription,
			Short:        "Verifies the resource file signature",
			SilenceUsage: true,
		},
	}

	f := v.cmd.PersistentFlags()
	f.StringVar(&v.name, "name", "", "name of the repository to pull")
	f.StringVar(&v.url, "url", "", "url of the repository to pull")
	f.StringVar(&v.resourceType, "type", "", "type of resource to pull")
	f.StringVar(&v.ignoreVersions, "ignore-versions", "", "versions to ignore while pulling")

	return v
}
