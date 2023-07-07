package fetcher_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/fetcher"
	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/fetcher/config"
	"gopkg.in/h2non/gock.v1"
)

func TestFetchContractFromRepository(t *testing.T) {
	t.Cleanup(gock.Off)

	repo := config.Repository{
		Name: "golang-task",
		URL:  "https://github.com/shortbrain/golang-tasks",
	}
	r := strings.TrimPrefix(repo.URL, "https://github.com/")

	gock.New("https://api.github.com").
		Get(fmt.Sprintf("repos/%s/releases", r)).
		Reply(200).
		File("testdata/releases.yaml")
	gock.New("https://github.com").
		Get(fmt.Sprintf("%s/releases/download/v1.0.0/contract.yaml", r)).
		Reply(200).
		File("config/testdata/contract.simple.yaml")

	client, err := api.DefaultRESTClient()
	if err != nil {
		t.Fatal(err)
	}
	m, err := fetcher.FetchContractsFromRepository(repo, client)
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 1 {
		t.Fatalf("Should have fetched only 1 version, fetched %d: %v", len(m), m)
	}
}
