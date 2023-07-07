package config_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/fetcher/config"
	"gopkg.in/h2non/gock.v1"
)

func TestLoadContractValid(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "contract.*.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		_, filename := filepath.Split(path)
		testname := filename[:len(filename)-len(filepath.Ext(path))]

		t.Run(testname, func(t *testing.T) {
			_, err := config.LoadContract(path)
			if err != nil {
				t.Fatalf("Shouldn't have errored out on %s : %v", path, err)
			}
		})
	}
}

func TestLoadContractInvalid(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "invalid*.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		_, filename := filepath.Split(path)
		testname := filename[:len(filename)-len(filepath.Ext(path))]

		t.Run(testname, func(t *testing.T) {
			_, err := config.LoadContract(path)
			if err == nil {
				t.Fatalf("Should have errored out on %s : %v", path, err)
			}
		})
	}
}

func TestLoadContractNonExisting(t *testing.T) {
	_, err := config.LoadContract("testdata/do-not-exists.yaml")
	if err == nil || !os.IsNotExist(errors.Unwrap(err)) {
		t.Fatalf("Should have errored out on non existing file : %v", err)
	}
}

func TestLoadContractFromURLValid(t *testing.T) {
	t.Cleanup(gock.Off)
	gock.New("https://foo.bar").
		Get("baz").
		Reply(200).
		File("testdata/contract.simple.yaml")
	c, err := config.LoadContractFromURL("https://foo.bar/baz")
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d: %v", len(c.Tasks), c)
	}
}

func TestLoadContractFromURLInvalid(t *testing.T) {
	t.Cleanup(gock.Off)
	gock.New("https://foo.bar").
		Get("baz").
		Reply(200).
		BodyString(`invalid content`)
	_, err := config.LoadContractFromURL("https://foo.bar/baz")
	if err == nil {
		t.Fatal("Should have errored out")
	}
}

func TestLoadContractFromURLHTTPError(t *testing.T) {
	t.Cleanup(gock.Off)
	gock.New("https://foo.bar").
		Get("baz").
		Reply(500).
		BodyString("wat")
	_, err := config.LoadContractFromURL("https://foo.bar/baz")
	if err == nil || err.Error() != "Status error: 500" {
		t.Fatalf("Should have errored out with status code 500: %v", err)
	}
}
