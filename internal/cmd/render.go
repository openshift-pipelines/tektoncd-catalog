package cmd

import (
	"fmt"
	"os"

	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/render"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/runner"

	"github.com/spf13/cobra"
)

type RenderCmd struct {
	cmd      *cobra.Command // cobra command definition
	resource string         // path to the resource file
}

var _ runner.SubCommand = &RenderCmd{}

const renderLongDescription = `# catalog-cd render

Renders the informed Tekton resource file as markdown, focusing on the most important attributes
which should always be part of the Task documentation.

The markdown generated contains the Workspaces, Params and Results formated as a mardown tables.
`

// Cmd shares the Cobra command instance.
func (r *RenderCmd) Cmd() *cobra.Command {
	return r.cmd
}

// Complete asserts a single argument is informed.
func (r *RenderCmd) Complete(_ *config.Config, args []string) error {
	if len(args) == 1 {
		r.resource = args[0]
		return nil
	}
	return fmt.Errorf("you must inform a single argument (%d)", len(args))
}

// Validate assert the informed resource file exists.
func (r *RenderCmd) Validate() error {
	_, err := os.Stat(r.resource)
	return err
}

// Run renders the resource as markdown.
func (r *RenderCmd) Run(cfg *config.Config) error {
	md, err := render.NewMarkdown(cfg, r.resource)
	if err != nil {
		return err
	}
	return md.Render()
}

// NewRenderCmd instantiate the "render" subcommand.
func NewRenderCmd() runner.SubCommand {
	return &RenderCmd{
		cmd: &cobra.Command{
			Use:   "render",
			Short: "Renders the informed Tekton resource file as markdown",
			Long:  renderLongDescription,
		},
	}
}
