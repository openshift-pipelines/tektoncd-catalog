package fetcher

import (
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/fetcher/config"
)

// TODO: prepare the "dest" workspace
//
//	(fetch the repository's `p` branch, …)
//
// TODO: fetch release assets
//   - fetch yamls
//   - fetch tests (yamls with kttl)
//   - fetch bundles, sbom, …
//
// TODO: extract in destination folder
// TODO: copy source/ to dest/ as well
//
//	(warn if there is conflicts)
//
// TODO: create a PR
func FetchContractsFromRepository(r config.Repository, client *api.RESTClient) (map[string]config.Contract, error) {
	m := map[string]config.Contract{}

	if !strings.HasPrefix(r.URL, "https://github.com") {
		return m, fmt.Errorf("Non-github repository not supported: %s", r.URL)
	}
	repo := strings.TrimPrefix(r.URL, "https://github.com/")
	versions, err := fetchVersions(repo, client)
	if err != nil {
		return m, fmt.Errorf("Failed to fetch versions from %s: %w", r.URL, err)
	}
	for _, v := range versions {
		if v.PreRelease || v.Draft {
			// Ignore drafts or pre-releases
			continue
		}
		var contractAsset Asset
		contractFound := false
		for _, a := range v.Assets {
			if a.Name == "contract.yaml" {
				contractFound = true
				contractAsset = a
				break
			}
		}
		if !contractFound {
			// FIXME(vdemeester) should we ignore or error out ?
			continue
		}
		// Load contract from asset
		contract, err := config.LoadContractFromURL(contractAsset.DownloadURL)
		if err != nil {
			return m, fmt.Errorf("Failed to load asset %s from %s: %w", contractAsset.Name, v.TagName, err)
		}
		m[v.TagName] = contract
	}
	return m, nil
}

func fetchVersions(github string, client *api.RESTClient) ([]Version, error) {
	versions := []Version{}
	err := client.Get(fmt.Sprintf("repos/%s/releases", github), &versions)
	if err != nil {
		return nil, err
	}
	return versions, nil
}

type Version struct {
	Name       string
	TagName    string `json:"tag_name"`
	Id         int
	Draft      bool
	PreRelease bool
	Assets     []Asset
	URL        string `json:"url"`
	TarballURL string `json:"tarball_url"`
}

type Asset struct {
	Id          int
	URL         string `json:"url"`
	Name        string
	Label       string
	ContentType string `json:"content_type"`
	State       string
	DownloadURL string `json:"browser_download_url"`
}
