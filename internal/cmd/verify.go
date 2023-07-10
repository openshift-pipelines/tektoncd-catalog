package cmd

import (
	"fmt"
	"os"

	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/runner"

	"github.com/sigstore/cosign/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/cmd/cosign/cli/verify"
	"github.com/spf13/cobra"
)

// VerifyCmd represents the "verify" subcommand to verify the signature of a resource file.
type VerifyCmd struct {
	cmd       *cobra.Command // cobra command definition
	key       string         // path to the public key file
	signature string         // path to the signature file
	blob      string         // path to the blob file
}

var _ runner.SubCommand = &VerifyCmd{}

const verifyLongDescription = `# catalog-cd verify

Verifies the cryptographic signature (--signature) of a arbitrary resource file (blob) using a public
key (--key) to assert the resource's author.

To sign a resource file, use "sigstore/cosign" as the follow example:

  $ cosign sign-blob \
      --key="/path/to/private.key" \
      --output-signature="/var/tmp/resource.sig" \
      /path/to/resource.yaml

The ".sig" file produced should be directly used as --signature flag. For example:

  $ catalog-cd verify \
      --key="/path/to/public-key.pub" \
      --signature="/var/tmp/resource.sig" \
      /path/to/resource.yaml
`

// Cmd exposes the cobra command instance.
func (v *VerifyCmd) Cmd() *cobra.Command {
	return v.cmd
}

// Complete asserts the required flags are informed, and the last argument is the resource file for
// signature verification.
func (v *VerifyCmd) Complete(_ *config.Config, args []string) error {
	if v.key == "" {
		return fmt.Errorf("flag --key is required")
	}
	if v.signature == "" {
		return fmt.Errorf("flag --signature is required")
	}

	if len(args) != 1 {
		return fmt.Errorf("you must inform a single blob file argument")
	}
	v.blob = args[0]
	return nil
}

// Validate asserts all the required files exists.
func (v *VerifyCmd) Validate() error {
	required := []string{
		v.key,
		v.signature,
		v.blob,
	}
	for _, f := range required {
		if _, err := os.Stat(f); err != nil {
			return err
		}
	}
	return nil
}

// Run wrapper around "cosign verify-blob" command.
func (v *VerifyCmd) Run(cfg *config.Config) error {
	cfg.Infof("# Public-Key: %q\n", v.key)
	cfg.Infof("#  Signature: %q\n", v.signature)
	cfg.Infof("#       Blob: %q\n", v.blob)
	return verify.VerifyBlobCmd(
		v.cmd.Context(),
		options.KeyOpts{KeyRef: v.key},
		"",
		"",
		"",
		"",
		"",
		v.signature,
		v.blob,
		"",
		"",
		"",
		"",
		"",
		false,
	)
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
		// informed via command-line flags
		key:       "",
		signature: "",
	}

	f := v.cmd.PersistentFlags()
	f.StringVar(&v.key, "key", v.key, "path to the public key file")
	f.StringVar(&v.signature, "signature", v.key, "path to the signature payload file")

	return v
}
