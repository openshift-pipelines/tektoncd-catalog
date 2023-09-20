package attestation

import (
	"context"

	"github.com/sigstore/cosign/v2/cmd/cosign/cli/generate"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/sign"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/verify"
)

// Attestation controls the sining and verification of resources.
type Attestation struct {
	rootOptions *options.RootOptions // general cosign settings
	keyOpts     options.KeyOpts      // public/private key reference

	privateKeyPass    []byte // stores the private key for signing
	base64            bool   // stores the signature using base64
	outputCertificate string // output certificate location
	tlogUpload        bool   // transaction log upload

	offline    bool // offline verification
	ignoreSCT  bool // ignore embedded SCT proof
	ignoreTlog bool // ignore transaction log
}

// GetPass prompts for the user private-key password only once, when the password is already
// stored it returns instead.
func (a *Attestation) GetPass(confirm bool) ([]byte, error) {
	if len(a.privateKeyPass) > 0 {
		return a.privateKeyPass, nil
	}

	var err error
	a.privateKeyPass, err = generate.GetPass(confirm)
	return a.privateKeyPass, err
}

// Sign signs the resource file (first argument) on the specificed signature location.
func (a *Attestation) Sign(payloadPath, outputSignature string) error {
	_, err := sign.SignBlobCmd(
		a.rootOptions,
		a.keyOpts,
		payloadPath,
		a.base64,
		outputSignature,
		a.outputCertificate,
		a.tlogUpload,
	)
	return err
}

// Verify verifies the resource signature
func (a *Attestation) Verify(ctx context.Context, blobRef, sigRef string) error {
	v := verify.VerifyBlobCmd{
		KeyOpts:    a.keyOpts,
		SigRef:     sigRef,
		IgnoreSCT:  a.ignoreSCT,
		IgnoreTlog: a.ignoreTlog,
		Offline:    a.offline,
	}
	return v.Exec(ctx, blobRef)
}

// NewAttestation instantiate the Attestation helper setting the default parameters expected
// for signing and verifying resources.
func NewAttestation(key string) (*Attestation, error) {
	o := &options.SignBlobOptions{}
	oidcClientSecret, err := o.OIDC.ClientSecret()
	if err != nil {
		return nil, err
	}

	keyOpts := options.KeyOpts{
		KeyRef:           key,
		Sk:               false,
		Slot:             "",
		FulcioURL:        options.DefaultFulcioURL,
		RekorURL:         options.DefaultRekorURL,
		OIDCIssuer:       options.DefaultOIDCIssuerURL,
		OIDCClientID:     "sigstore",
		OIDCClientSecret: oidcClientSecret,
		SkipConfirmation: true,
	}

	a := &Attestation{
		rootOptions:       &options.RootOptions{},
		keyOpts:           keyOpts,
		base64:            true,
		ignoreSCT:         false,
		ignoreTlog:        true,
		offline:           true,
		outputCertificate: "",
		tlogUpload:        false,
	}
	a.keyOpts.PassFunc = a.GetPass

	return a, nil
}
