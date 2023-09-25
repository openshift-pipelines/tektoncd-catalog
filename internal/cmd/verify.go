package cmd

import (
	"context"
	"fmt"

	"github.com/openshift-pipelines/tektoncd-catalog/internal/attestation"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/contract"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/runner"

	"github.com/spf13/cobra"
)

// VerifyCmd represents the "verify" subcommand to verify the signature of a resource file.
type VerifyCmd struct {
	cmd *cobra.Command // cobra command definition
	c   *contract.Contract

	publicKey string // path to the public key file
}

var _ runner.SubCommand = &VerifyCmd{}

const verifyLongDescription = `# catalog-cd verify

Verifies the signature of all resources described on the contract. The subcommand takes
either a contract file as argument, or a directory containing the contract using default
name. By default it searches the current directory.

In order to verify the signature the public-key is required, it's specified either on the the
catalog contract, or using the flag "--public-key".
`

// Cmd exposes the cobra command instance.
func (v *VerifyCmd) Cmd() *cobra.Command {
	return v.cmd
}

// Complete asserts the required flags are informed, and the last argument is the resource
// file for signature verification.
func (v *VerifyCmd) Complete(_ *config.Config, args []string) error {
	var err error
	v.c, err = LoadContractFromArgs(args)
	return err
}

// Validate asserts all the required files exists.
func (v *VerifyCmd) Validate() error {
	var err error
	if v.publicKey == "" {
		v.publicKey, err = v.c.GetPublicKey()
	}
	return err
}

// Run wrapper around "cosign verify-blob" command.
func (v *VerifyCmd) Run(cfg *config.Config) error {
	cfg.Infof("# Public-Key: %q\n", v.publicKey)

	helper, err := attestation.NewAttestation(v.publicKey)
	if err != nil {
		return err
	}
	return v.c.VerifyResources(v.cmd.Context(), func(ctx context.Context, blobRef, sigRef string) error {
		fmt.Printf("# Verifying resource %q against signature %q...\n", blobRef, sigRef)
		return helper.Verify(ctx, blobRef, sigRef)
	})
}

// NewVerifyCmd instantiates the "verify" subcommand.
func NewVerifyCmd() runner.SubCommand {
	v := &VerifyCmd{
		cmd: &cobra.Command{
			Use:          "verify",
			Args:         cobra.ExactArgs(1),
			Long:         verifyLongDescription,
			Short:        "Verifies the resource file signature",
			SilenceUsage: true,
		},
	}

	f := v.cmd.PersistentFlags()
	f.StringVar(&v.publicKey, "public-key", "", "path to the public key file")

	return v
}
