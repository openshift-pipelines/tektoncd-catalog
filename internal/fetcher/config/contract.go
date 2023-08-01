package config

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"sigs.k8s.io/yaml"
)

// VersionnedContract is a "fake" struct to parse a contract and get its version number
type VersionnedContract struct {
	Version string
}

// Contract is a representation of the configuration/contract that is attached to releases
// TODO(vdemeester) Rename this to ContractV0 when adding new versions
type Contract struct {
	Tasks     []Task
	Pipelines []Pipeline
}

type Task struct {
	Name   string
	File   string
	Bundle string
	// TODO: Tests
	// TODO: SBOM
	// TODO: Attestation
	// TODO: Signatures
}

type Pipeline struct {
	Name   string
	File   string
	Bundle string
	// TODO: try this out — linking Pipeline and Tasks
	Tasks []Task
}

func LoadContract(filename string) (Contract, error) {
	var vc VersionnedContract
	var c Contract
	data, err := os.ReadFile(filename)
	if err != nil {
		return c, fmt.Errorf("Could not load contract from %s: %w", filename, err)
	}
	if err := yaml.Unmarshal(data, &vc); err != nil {
		return c, fmt.Errorf("Could not load contract version from %s: %w", filename, err)
	}
	// FIXME(vdemeester) change this once we support multiple version
	if vc.Version != "0" {
		return c, fmt.Errorf("Contract version %s is not supported", vc.Version)
	}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("Could not load contract from %s: %w", filename, err)
	}
	return c, nil
}

func LoadContractFromURL(url string) (Contract, error) {
	var c Contract
	var vc VersionnedContract

	resp, err := http.Get(url)
	if err != nil {
		return c, fmt.Errorf("Could not load contract from %s: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return c, fmt.Errorf("Status error: %v", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c, fmt.Errorf("Could not load contract from %s: %w", url, err)
	}
	if err := yaml.Unmarshal(data, &vc); err != nil {
		return c, fmt.Errorf("Could not load contract version from %s: %w", url, err)
	}
	if vc.Version != "0" {
		return c, fmt.Errorf("Contract version %s is not supported", vc.Version)
	}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("Could not load contract from %s: %w", url, err)
	}
	return c, nil
}
