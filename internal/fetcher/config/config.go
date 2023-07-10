// Package config holds the "configuration" that is used for fetching information from external repositories.
package config

import (
	"fmt"
	"os"

	"sigs.k8s.io/yaml"
)

// External is a representation of the configuration for specifying repositories we have to pull from
type External struct {
	// Repositories defines the repositories to pull from
	Repositories []Repository
}

// Repository represent a git repository
type Repository struct {
	Name string
	URL  string
	// Type defines the type to fetch (Task, Pipeline, …)
	Types          []string
	IgnoreVersions []string `json:"ignore-versions"`
}

func LoadExternal(filename string) (External, error) {
	var c External
	data, err := os.ReadFile(filename)
	if err != nil {
		return c, fmt.Errorf("Could not load external configuration from %s: %w", filename, err)
	}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("Could not load external configuration from %s: %w", filename, err)
	}
	return c, nil
}
