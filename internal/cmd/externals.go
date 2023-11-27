package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
	fc "github.com/openshift-pipelines/tektoncd-catalog/internal/fetcher/config"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/runner"
	"github.com/spf13/cobra"
)

// ExternalsCmd represents the "externals" subcommand to externals the signature of a resource file.
type ExternalsCmd struct {
	cmd    *cobra.Command // cobra command definition
	config string         // path for the catalog configuration file
}

var _ runner.SubCommand = &ExternalsCmd{}

const externalsLongDescription = `# catalog-cd externals

TODO
`

// Cmd exposes the cobra command instance.
func (v *ExternalsCmd) Cmd() *cobra.Command {
	return v.cmd
}

// Complete asserts the required flags are informed, and the last argument is the resource file for
// signature verification.
func (v *ExternalsCmd) Complete(_ *config.Config, args []string) error {
	if v.config == "" {
		return fmt.Errorf("flag --config is required")
	}

	if len(args) != 0 {
		return fmt.Errorf("externals takes no argument")
	}
	return nil
}

// Validate asserts all the required files exists.
func (v *ExternalsCmd) Validate() error {
	required := []string{
		v.config,
	}
	for _, f := range required {
		if _, err := os.Stat(f); err != nil {
			return err
		}
	}
	return nil
}

type GitHubRunObject struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	Type           string `json:"type"`
	IgnoreVersions string `json:"ignoreVersions"`
}

type GitHubMatrixObject struct {
	Include []GitHubRunObject `json:"include"`
}

// Run
func (v *ExternalsCmd) Run(cfg *config.Config) error {
	e, err := fc.LoadExternal(v.config)
	if err != nil {
		return err
	}
	m := GitHubMatrixObject{}
	for _, repository := range e.Repositories {
		types := repository.Types
		if len(types) == 0 {
			types = []string{"tasks", "pipelines"}
		}
		ignoreVersions := ""
		if len(repository.IgnoreVersions) > 0 {
			ignoreVersions = strings.Join(repository.IgnoreVersions, ",")
		}
		for _, t := range types {
			name := repository.Name
			if name == "" {
				name = path.Base(repository.URL)
			}
			o := GitHubRunObject{
				Name:           name,
				URL:            repository.URL,
				Type:           t,
				IgnoreVersions: ignoreVersions,
			}
			m.Include = append(m.Include, o)
		}
	}
	j, err := json.Marshal(m)
	if err != nil {
		return err
	}
	fmt.Fprintf(cfg.Stream.Out, "%s\n", j)
	return nil
}

// NewCatalogExternalsCmd instantiates the "externals" subcommand.
func NewCatalogExternalsCmd() runner.SubCommand {
	v := &ExternalsCmd{
		cmd: &cobra.Command{
			Use:          "externals",
			Args:         cobra.ExactArgs(0),
			Long:         externalsLongDescription,
			Short:        "Verifies the resource file signature",
			SilenceUsage: true,
		},
	}

	f := v.cmd.PersistentFlags()
	f.StringVar(&v.config, "config", v.config, "path of the catalog configuration file")

	return v
}
