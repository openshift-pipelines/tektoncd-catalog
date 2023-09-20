package contract

import (
	"context"
	"fmt"

	"github.com/openshift-pipelines/tektoncd-catalog/internal/resource"
)

// TektonResource contains a Tekton resource reference, as in a Task or Pipeline.
type TektonResource struct {
	// Name Tekton resource name, the Task or Pipeline actual name.
	Name string `json:"name"`
	// Version Tekton resource version.
	Version string `json:"version"`
	// Filename starting from the repository root, the relative path to the resource file.
	Filename string `json:"filename"`
	// Checksum ".filename"'s SHA256 sum, validates resource payload after network transfer.
	Checksum string `json:"checksum"`
	// Signature Tekton resource signature, either the signature payload, or relative
	// location to the signature file. By default, it uses the ".filename" attributed
	// followed by ".sig" extension.
	Signature string `json:"signature"`
}

// Resources inventory of all Tekton resources managed by the repository.
type Resources struct {
	// Tasks List of Tekton Tasks.
	Tasks []*TektonResource `json:"tasks"`
	// Pipelines List of Tekton Pipelines.
	Pipelines []*TektonResource `json:"pipelines"`
}

// ResourceSignFn function to perform the resource (file) signature. Parameters:
//   - resource-file: resource file location to be signed
//   - signature-file: where the signature file should be stored
type ResourceSignFn func(_, _ string) error

// ResourceVerifySignatureFn function to perform the signature verification. Parameters:
//   - context: shared context
//   - resource-file: the resource file
//   - signature-file: the respective signature file
type ResourceVerifySignatureFn func(_ context.Context, _, _ string) error

// SignResources runs the informed function against each catalog resource, the expected
// signature file created is updated on "this" contract instance.
func (c *Contract) SignResources(fn ResourceSignFn) error {
	for _, r := range append(c.Catalog.Resources.Tasks, c.Catalog.Resources.Pipelines...) {
		signatureFile := fmt.Sprintf("%s.%s", r.Filename, SignatureExtension)
		if err := fn(r.Filename, signatureFile); err != nil {
			return err
		}
		r.Signature = signatureFile
	}
	return nil
}

// VerifyResources runs the informed function against each catalog resource, when error is
// returned the signature verification process fail.
func (c *Contract) VerifyResources(ctx context.Context, fn ResourceVerifySignatureFn) error {
	for _, r := range append(c.Catalog.Resources.Tasks, c.Catalog.Resources.Pipelines...) {
		if err := fn(ctx, r.Filename, r.Signature); err != nil {
			return err
		}
	}
	return nil
}

// AddResourceFile adds a resource file on the contract, making sure it's a Tekton resource
// file and uses the "kind" to guide on which attribute the resource will be appended.
func (c *Contract) AddResourceFile(resourceFile string, version string) error {
	// parsing the resource as a kubernetes unstructured type to read it's name and kind
	u, err := resource.ReadAndDecodeResourceFile(resourceFile)
	if err != nil {
		return err
	}
	// making sure it's a tekton kubernetes resource, on the supported versions
	if err = isResourceSupported(u); err != nil {
		return err
	}

	sha256sum, err := CalculateSHA256Sum(resourceFile)
	if err != nil {
		return err
	}

	tr := TektonResource{
		Name:     u.GetName(),
		Version:  version,
		Filename: resourceFile,
		Checksum: sha256sum,
	}

	switch kind := u.GetKind(); kind {
	case "Task":
		c.Catalog.Resources.Tasks = append(c.Catalog.Resources.Tasks, &tr)
	case "Pipeline":
		c.Catalog.Resources.Pipelines = append(c.Catalog.Resources.Pipelines, &tr)
	default:
		return fmt.Errorf("%w: resource kind %q", ErrTektonResourceUnsupported, kind)
	}
	return nil
}
