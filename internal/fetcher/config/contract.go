package config

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"sigs.k8s.io/yaml"
)

// TODO(vdemeest) version contracts

// Contract is a representation of the configuration/contract that is attached to releases
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
	var c Contract
	data, err := os.ReadFile(filename)
	if err != nil {
		return c, fmt.Errorf("Could not load contract from %s: %w", filename, err)
	}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("Could not load contract from %s: %w", filename, err)
	}
	return c, nil
}

func LoadContractFromURL(url string) (Contract, error) {
	var c Contract

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
	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("Could not load contract from %s: %w", url, err)
	}
	return c, nil
}
