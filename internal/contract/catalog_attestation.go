package contract

import (
	"fmt"

	"github.com/go-errors/errors"
)

// ErrAttestationPublicKeyEmpty marks the public-key is not yet set.
var ErrAttestationPublicKeyEmpty = errors.New("public-key is empty")

// Attestation holds the attributes needed for the software supply chain security.
type Attestation struct {
	// PublicKey path to the public key file, KMS URI or Kubernetes Secret.
	PublicKey string `json:"publicKey"`
}

// GetPublicKey accessor to the attestation's public-key, emits error when not set.
func (c *Contract) GetPublicKey() (string, error) {
	if c.Catalog.Attestation == nil || c.Catalog.Attestation.PublicKey == "" {
		return "", fmt.Errorf("%w: .catalog.attestation is not set", ErrAttestationPublicKeyEmpty)
	}
	return c.Catalog.Attestation.PublicKey, nil
}
