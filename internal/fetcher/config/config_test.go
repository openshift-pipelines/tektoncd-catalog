package config_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/openshift-pipelines/tektoncd-catalog/internal/fetcher/config"
)

func TestLoadExternalValid(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "external.*.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		_, filename := filepath.Split(path)
		testname := filename[:len(filename)-len(filepath.Ext(path))]

		t.Run(testname, func(t *testing.T) {
			_, err := config.LoadExternal(path)
			if err != nil {
				t.Fatalf("Shouldn't have errored out on %s : %v", path, err)
			}
		})
	}
}

func TestLoadExternalInvalid(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "invalid*.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		_, filename := filepath.Split(path)
		testname := filename[:len(filename)-len(filepath.Ext(path))]

		t.Run(testname, func(t *testing.T) {
			_, err := config.LoadExternal(path)
			if err == nil {
				t.Fatalf("Should have errored out on %s : %v", path, err)
			}
		})
	}
}

func TestLoadExternalNonExisting(t *testing.T) {
	_, err := config.LoadExternal("testdata/do-not-exists.yaml")
	if err == nil || !os.IsNotExist(errors.Unwrap(err)) {
		t.Fatalf("Should have errored out on non existing file : %v", err)
	}
}
