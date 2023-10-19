package catalog

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/fetcher"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/fetcher/config"
)

func FetchFromExternals(e config.External, client *api.RESTClient) (Catalog, error) {
	c := Catalog{
		Tasks:     map[string]Task{},
		Pipelines: map[string]Pipeline{},
	}
	for _, r := range e.Repositories {
		fmt.Fprintln(os.Stderr, "Fetching", r.Name, "("+r.URL+")")
		var fetchTask, fetchPipeline bool
		if r.Types == nil {
			fetchTask = true
			fetchPipeline = true
		} else {
			for _, t := range r.Types {
				if t == "tasks" {
					fetchTask = true
				}
				if t == "pipelines" {
					fetchPipeline = true
				}
			}
		}
		m, err := fetcher.FetchContractsFromRepository(r, client)
		if err != nil {
			return c, err
		}
		for _, v := range r.IgnoreVersions {
			if _, ok := m[v]; ok {
				// Remove ignored versions from map
				delete(m, v)
			}
		}

		for version, contract := range m {
			if fetchTask {
				if contract.Catalog.Resources != nil {
					for _, task := range contract.Catalog.Resources.Tasks {
						if _, ok := c.Tasks[task.Name]; !ok {
							// task doesn't exists yet, creating it
							c.Tasks[task.Name] = Task{
								Versions: map[string]VersionnedTask{},
							}
						}
						if _, ok := c.Tasks[task.Name].Versions[version]; ok {
							//  name/version confict
							return c, fmt.Errorf("Task %s has a version conflict (%s)", task.Name, r.URL)
						}
						downloadURL := task.Filename
						if !strings.HasPrefix(task.Filename, "https://") {
							downloadURL = fmt.Sprintf("%s/releases/download/%s/%s", r.URL, version, task.Filename)
						}
						c.Tasks[task.Name].Versions[version] = VersionnedTask{
							DownloadURL: downloadURL,
							// Bundle:      task.Bundle, // FIXME: add bundle support
						}
					}
				}
			}
			if fetchPipeline {
				if contract.Catalog.Resources != nil {
					for _, pipeline := range contract.Catalog.Resources.Pipelines {
						if _, ok := c.Pipelines[pipeline.Name]; !ok {
							// pipeline doesn't exists yet, creating it
							c.Pipelines[pipeline.Name] = Pipeline{
								Versions: map[string]VersionnedPipeline{},
							}
						}
						if _, ok := c.Pipelines[pipeline.Name].Versions[version]; ok {
							// name/version confict
							return c, fmt.Errorf("Pipeline %s has a version conflict (%s)", pipeline.Name, r.URL)
						}
						downloadURL := pipeline.Filename
						if !strings.HasPrefix(pipeline.Filename, "https://") {
							downloadURL = fmt.Sprintf("%s/releases/download/%s/%s", r.URL, version, pipeline.Filename)
						}
						c.Pipelines[pipeline.Name].Versions[version] = VersionnedPipeline{
							DownloadURL: downloadURL,
							// Bundle:      pipeline.Bundle, // FIXME: add bundle support
						}
					}
				}
			}
		}
	}
	return c, nil
}

func GenerateFilesystem(path string, c Catalog) error {
	if err := generateTasksFilesystem(filepath.Join(path, "tasks"), c.Tasks); err != nil {
		return fmt.Errorf("Failed to create the tasks filesystem: %w", err)
	}
	if err := generatePipelinesFilesystem(filepath.Join(path, "pipelines"), c.Pipelines); err != nil {
		return fmt.Errorf("Failed to create the tasks filesystem: %w", err)
	}
	return nil
}

func generateTasksFilesystem(path string, tasks map[string]Task) error {
	for name, t := range tasks {
		for version, task := range t.Versions {
			taskfolder := filepath.Join(path, name, version)
			if err := os.MkdirAll(taskfolder, os.ModePerm); err != nil {
				return err
			}
			taskfile := filepath.Join(taskfolder, fmt.Sprintf("%s.yaml", name))
			if err := fetchAndWrite(taskfile, task.DownloadURL); err != nil {
				return fmt.Errorf("Couldn't fetch %s in %s: %w", task.DownloadURL, taskfile, err)
			}
		}
	}
	return nil
}

func generatePipelinesFilesystem(path string, pipelines map[string]Pipeline) error {
	for name, t := range pipelines {
		for version, pipeline := range t.Versions {
			pipelinefolder := filepath.Join(path, name, version)
			if err := os.MkdirAll(filepath.Join(path, name, version), os.ModePerm); err != nil {
				return err
			}
			pipelinefile := filepath.Join(pipelinefolder, fmt.Sprintf("%s.yaml", name))
			if err := fetchAndWrite(pipelinefile, pipeline.DownloadURL); err != nil {
				return fmt.Errorf("Couldn't fetch %s in %s: %w", pipeline.DownloadURL, pipelinefile, err)
			}
		}
	}
	return nil
}

func fetchAndWrite(file, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status error: %v", resp.StatusCode)
	}
	w, err := os.Create(file)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// Catalog is a struct that represent a "file-based" catalog
// file-based catalog
type Catalog struct {
	Tasks     map[string]Task
	Pipelines map[string]Pipeline
}

type Task struct {
	Versions map[string]VersionnedTask
}

type VersionnedTask struct {
	DownloadURL string
	Bundle      string
}

type Pipeline struct {
	Versions map[string]VersionnedPipeline
}

type VersionnedPipeline struct {
	DownloadURL string
	Bundle      string
}
