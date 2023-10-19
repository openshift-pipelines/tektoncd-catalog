package catalog_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/catalog"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/fetcher/config"
	"gopkg.in/h2non/gock.v1"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/fs"
	"gotest.tools/v3/golden"
)

func TestFetchFromExternal(t *testing.T) {
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
		File("testdata/contract.simple.yaml")

	client, err := api.DefaultRESTClient()
	if err != nil {
		t.Fatal(err)
	}
	e := config.External{
		Repositories: []config.Repository{{
			Name:  "sbr-golang",
			URL:   "https://github.com/shortbrain/golang-tasks",
			Types: []string{"tasks"},
		}},
	}
	c, err := catalog.FetchFromExternals(e, client)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Tasks) != 2 {
		t.Fatalf("Should have created a catalog with only 2 task, got %d: %v", len(c.Tasks), c.Tasks)
	}
}

func TestGenerateFilesystem(t *testing.T) {
	t.Cleanup(gock.Off)

	gock.New("https://fake.host").
		Get("git-clone-0.1.0.yaml").
		Reply(200).
		File("testdata/git-clone.yaml")
	gock.New("https://fake.host").
		Get("git-clone-1.1.0.yaml").
		Reply(200).
		File("testdata/git-clone.yaml")
	gock.New("https://fake.host").
		Get("golang-build-0.2.0.yaml").
		Reply(200).
		File("testdata/golang-build.yaml")
	gock.New("https://fake.host").
		Get("my-pipeline-1.0.0.yaml").
		Reply(200).
		File("testdata/my-pipeline.yaml")

	dir := fs.NewDir(t, "catalog")
	defer dir.Remove()

	c := catalog.Catalog{
		Tasks: map[string]catalog.Task{
			"git-clone": {
				Versions: map[string]catalog.VersionnedTask{
					"0.1.0": {
						DownloadURL: "https://fake.host/git-clone-0.1.0.yaml",
						Bundle:      "fake.host/gitclone:0.1.0",
					},
					"1.1.0": {
						DownloadURL: "https://fake.host/git-clone-1.1.0.yaml",
						Bundle:      "fake.host/gitclone:1.1.0",
					},
				},
			},
			"golang-build": {
				Versions: map[string]catalog.VersionnedTask{
					"0.2.0": {
						DownloadURL: "https://fake.host/golang-build-0.2.0.yaml",
					},
				},
			},
		},
		Pipelines: map[string]catalog.Pipeline{
			"my-pipeline": {
				Versions: map[string]catalog.VersionnedPipeline{
					"1.0.0": {
						DownloadURL: "https://fake.host/my-pipeline-1.0.0.yaml",
					},
				},
			},
		},
	}
	err := catalog.GenerateFilesystem(dir.Path(), c)
	if err != nil {
		t.Fatal(err)
	}
	expected := fs.Expected(t,
		fs.WithDir("tasks",
			fs.WithDir("git-clone",
				fs.WithDir("0.1.0", fs.WithFile("git-clone.yaml", "", fs.WithBytes(golden.Get(t, "git-clone.yaml")))),
				fs.WithDir("1.1.0", fs.WithFile("git-clone.yaml", "", fs.WithBytes(golden.Get(t, "git-clone.yaml")))),
			),
			fs.WithDir("golang-build",
				fs.WithDir("0.2.0", fs.WithFile("golang-build.yaml", "", fs.WithBytes(golden.Get(t, "golang-build.yaml")))),
			),
		),
		fs.WithDir("pipelines",
			fs.WithDir("my-pipeline",
				fs.WithDir("1.0.0", fs.WithFile("my-pipeline.yaml", "", fs.WithBytes(golden.Get(t, "my-pipeline.yaml")))),
			),
		),
	)

	assert.Assert(t, fs.Equal(dir.Path(), expected))
}
