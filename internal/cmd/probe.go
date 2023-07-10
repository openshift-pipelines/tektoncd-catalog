package cmd

import (
	"fmt"
	"regexp"
	"time"

	"github.com/openshift-pipelines/tektoncd-catalog/internal/assert"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/flags"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/probe"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/runner"

	"github.com/spf13/cobra"
)

type ProbeCmd struct {
	cmd           *cobra.Command  // cobra command definition
	kind          string          // kind of tekton resource being probed
	rawParams     []string        // input params
	rawWorkspaces []string        // input workspaces
	assertResults []regexp.Regexp // slice of regexps to assert output
	name          string          // name of the tekton resource to be tested
}

var _ runner.SubCommand = &ProbeCmd{}

func (p *ProbeCmd) Cmd() *cobra.Command {
	return p.cmd
}

func (p *ProbeCmd) Complete(_ *config.Config, args []string) error {
	if len(args) == 1 {
		p.name = args[0]
		return nil
	}
	return fmt.Errorf("you must inform a single argument (%d)", len(args))
}

func (p *ProbeCmd) Validate() error {
	return nil
}

func (p *ProbeCmd) Run(cfg *config.Config) error {
	resource, err := probe.NewProbe(cfg, p.kind, p.rawWorkspaces, p.rawParams)
	if err != nil {
		return err
	}
	subject, err := resource.Run(p.cmd.Context(), p.name)
	if err != nil {
		return err
	}

	// waiting a few seconds before continue, waiting for possible task status synchronization
	time.Sleep(5 * time.Second)

	expect, err := assert.NewAssert(p.cmd.Context(), cfg, subject)
	if err != nil {
		return err
	}
	if err = expect.Status(); err != nil {
		return err
	}
	return expect.Results(p.assertResults)
}

func NewProbeCmd() runner.SubCommand {
	p := &ProbeCmd{
		cmd: &cobra.Command{
			Use:          "probe",
			Short:        "TODO",
			Long:         "TODO",
			Args:         cobra.MinimumNArgs(1),
			SilenceUsage: true,
		},
		kind:          "task",
		assertResults: []regexp.Regexp{},
	}

	f := p.cmd.PersistentFlags()

	f.StringVar(&p.kind, "kind", p.kind,
		"Kind of Tekton Resource being probed")
	f.StringArrayVar(&p.rawParams, "param", []string{},
		"Param key-value, split by equal sign")
	f.StringArrayVar(&p.rawWorkspaces, "workspace", []string{},
		"Workspace key-value, split by equal sign")
	f.Var(flags.NewRegexpValue(&p.assertResults), "assert-result",
		"Regular expression to assert the output produced")

	return p
}
