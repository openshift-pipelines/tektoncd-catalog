package linter

import (
	"errors"
	// "fmt"

	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/resource"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Linter enforce linting rules on the informed Tekton resource file.
type Linter struct {
	cfg *config.Config             // global configuration
	u   *unstructured.Unstructured // object instance
}

var (
	ErrInvalidWorkspace = errors.New("invalid workspace definition")
	ErrInvalidParam     = errors.New("invalid param definition")
	ErrInvalidResult    = errors.New("invalid results definition")
)

// workspaces extract and lint the workspaces.
func (l *Linter) workspaces() error {
	l.cfg.Infof("# Inspecting workspaces...\n")
	wbs, err := GetNestedSlice(l.u, "spec", "workspaces")
	if err != nil {
		return err
	}
	if err = lowercaseSliceMapLinter(wbs); err != nil {
		// return fmt.Errorf("%w: %w", ErrInvalidWorkspace, err)
		l.cfg.Errorf("%s: %s", ErrInvalidWorkspace, err)
		return nil
	}
	l.cfg.Infof("# Workspaces are following best practices!\n")
	return nil
}

// params extract and lint the params.
func (l *Linter) params() error {
	l.cfg.Infof("# Inspecting params...\n")
	ps, err := GetNestedSlice(l.u, "spec", "params")
	if err != nil {
		return err
	}
	if err = uppercaseSliceMapLinter(ps); err != nil {
		// return fmt.Errorf("%w: %w", ErrInvalidParam, err)
		l.cfg.Errorf("%s: %s", ErrInvalidParam, err)
		return nil
	}
	l.cfg.Infof("# Params are following best practices!\n")
	return nil
}

// results extract and lint the results.
func (l *Linter) results() error {
	l.cfg.Infof("# Inspecting results...\n")
	rs, err := GetNestedSlice(l.u, "spec", "results")
	if err != nil {
		return err
	}
	if err = uppercaseSliceMapLinter(rs); err != nil {
		// return fmt.Errorf("%w: %w", ErrInvalidResult, err)
		l.cfg.Errorf("%s: %s", ErrInvalidResult, err)
		return nil
	}
	l.cfg.Infof("# Results are following best practices!\n")
	return nil
}

// Enforce lint workspaces, params and results.
func (l *Linter) Enforce() error {
	err := l.workspaces()
	if err != nil {
		return err
	}
	if err = l.params(); err != nil {
		return err
	}
	return l.results()
}

// NewLinter instantiate the resource linter by reading and decoding the resource file.
func NewLinter(cfg *config.Config, resourceFile string) (*Linter, error) {
	cfg.Infof("# Linting resource file %q...\n", resourceFile)
	u, err := resource.ReadAndDecodeResourceFile(resourceFile)
	if err != nil {
		return nil, err
	}
	cfg.Infof("# Name=%q, APIVersion=%q, Kind=%q\n", u.GetName(), u.GetAPIVersion(), u.GetKind())

	return &Linter{
		cfg: cfg,
		u:   u,
	}, nil
}
