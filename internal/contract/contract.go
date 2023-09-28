package contract

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"gopkg.in/yaml.v3"
)

const (
	// Version current contract version.
	Version = "v1"
	// Filename default contract file name.
	Filename = ".catalog.yaml"
	// SignatureExtension
	SignatureExtension = "sig"
)

// Repository contains the general repository information, including metadata to categorize
// and describe the repository contents, objective, ecosystem, etc.
type Repository struct {
	// Description long description text.
	Description string `json:"description"`
}

// ResourceProbe describes a single test-case for a Tekton resource managed by the
// repository, serves as inputs for "catalog-cd probe".
type ResourceProbe struct {
	// Name testa-case unique name.
	Name string `json:"name"`
	// ResourceName the name of the Tekton resource, present on ".catalog.resources".
	ResourceName string `json:"resourceName"`
	// Workspaces slice of Tekton workspace-bindings for the test-case.
	Workspaces []v1beta1.WorkspaceBinding `json:"workspaces"`
	// Params slice of Tekton Params for the test-case
	Params []v1beta1.Param `json:"params"`
}

// Probe contains all the test-cases for the Tekton resources managed by the repository.
type Probe struct {
	// Tasks Tekton Tasks tests.
	Tasks []ResourceProbe `json:"tasks"`
	// Pipelines Tekton Pipelines tests.
	Pipelines []ResourceProbe `json:"pipelines"`
}

// Catalog describes the contents of a repository part of a "catalog" of Tekton resources,
// including repository metadata, inventory of Tekton resources, test-cases and more.
type Catalog struct {
	Repository  *Repository  `json:"repository"`  // repository long description
	Attestation *Attestation `json:"attestation"` // software supply provenance
	Resources   *Resources   `json:"resources"`   // inventory of Tekton resources
	Probe       *Probe       `json:"probe"`       // test-cases for the managed resources
}

// Contract contains a versioned catalog.
type Contract struct {
	file    string  // contract file full path
	Version string  `json:"version"` // contract version
	Catalog Catalog `json:"catalog"` // tekton resources catalog
}

// Print renders the YAML representation of the current contract.
func (c *Contract) Print() ([]byte, error) {
	var b bytes.Buffer
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)
	if err := enc.Encode(c); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Save saves the contract on the original file.
func (c *Contract) Save() error {
	if c.file == "" {
		return fmt.Errorf("contract file location is not set")
	}
	return c.SaveAs(c.file)
}

// SaveAs writes itself on the informed file path.
func (c *Contract) SaveAs(file string) error {
	payload, err := c.Print()
	if err != nil {
		return err
	}
	return os.WriteFile(file, payload, 0644)
}

// NewContractEmpty instantiates a new Contract{} with empty attributes.
func NewContractEmpty() *Contract {
	return &Contract{
		Version: Version,
		Catalog: Catalog{
			Repository:  &Repository{},
			Attestation: &Attestation{},
			Resources: &Resources{
				Tasks:     []*TektonResource{},
				Pipelines: []*TektonResource{},
			},
			Probe: &Probe{
				Tasks:     []ResourceProbe{},
				Pipelines: []ResourceProbe{},
			},
		},
	}
}

// NewContractFromFile instantiates a new Contract{} from a YAML file.
func NewContractFromFile(location string) (*Contract, error) {
	// contract yaml file location
	var file string

	// when the location is a directory, it assumes the directory contains a default catalog
	// file name inside, otherwise the location is assumed to be the actual file
	info, _ := os.Stat(location)
	if info.IsDir() {
		file = path.Join(location, Filename)
	} else {
		file = location
	}

	payload, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return NewContractFromData(payload)
}

// NewContractFromURL instantiates a new Contract{} from a URL.
func NewContractFromURL(url string) (*Contract, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Could not load contract from %s: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not load contract from %s: %w", url, err)
	}
	return NewContractFromData(data)
}

// NewContractFromData instantiates a new Contract{} from a YAML payload.
func NewContractFromData(payload []byte) (*Contract, error) {
	c := Contract{}
	if err := yaml.Unmarshal(payload, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
