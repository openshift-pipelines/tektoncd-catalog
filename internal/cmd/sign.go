package cmd

import (
	"fmt"

	"github.com/openshift-pipelines/tektoncd-catalog/internal/attestation"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/contract"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/runner"

	"github.com/spf13/cobra"
)

// SignCmd subcommand "sign" to handles signing contract resources.
type SignCmd struct {
	cmd *cobra.Command     // cobra command definition
	c   *contract.Contract // catalog contract instance

	privateKey string // private key location
}

var _ runner.SubCommand = &SignCmd{}

const signLongDescription = `# catalog-cd sign

Sign the catalog contract resources on the informed directory, or catalog file. By default it
assumes the current directory.

To sign the resources the subcommand requires a private-key ("--private-key" flag), and may
ask for the password when trying to interact with a encripted key.
`

// Cmd exposes the cobra command instance.
func (s *SignCmd) Cmd() *cobra.Command {
	return s.cmd
}

// Complete loads the contract file from the location informed on the first argument.
func (s *SignCmd) Complete(_ *config.Config, args []string) error {
	var err error
	s.c, err = LoadContractFromArgs(args)
	return err
}

// Validate implements runner.SubCommand.
func (*SignCmd) Validate() error {
	return nil
}

// Run perform the resource signing.
func (s *SignCmd) Run(_ *config.Config) error {
	helper, err := attestation.NewAttestation(s.privateKey)
	if err != nil {
		return err
	}
	if err = s.c.SignResources(func(payladPath, outputSignature string) error {
		fmt.Printf("# Signing resource %q on %q...\n", payladPath, outputSignature)
		return helper.Sign(payladPath, outputSignature)
	}); err != nil {
		return err
	}
	return s.c.Save()
}

// NewSignCmd instantiate the SignCmd and flags.
func NewSignCmd() runner.SubCommand {
	s := &SignCmd{
		cmd: &cobra.Command{
			Use:          "sign [flags]",
			Short:        "Signs Tekton Pipelines resources",
			Long:         signLongDescription,
			Args:         cobra.MaximumNArgs(1),
			SilenceUsage: true,
		},
	}

	f := s.cmd.PersistentFlags()
	f.StringVar(&s.privateKey, "private-key", "", "private key file location")

	if err := s.cmd.MarkPersistentFlagRequired("private-key"); err != nil {
		panic(err)
	}

	return s
}
