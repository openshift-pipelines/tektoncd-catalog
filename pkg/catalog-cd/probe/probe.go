package probe

import (
	"context"
	"errors"
	"fmt"

	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/config"

	"github.com/tektoncd/cli/pkg/workspaces"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// Probe able to instantiate and run a Tekton resource, which will be subject to assertion later on.
type Probe interface {
	// Run runs the resource until completion.
	Run(_ context.Context, _ string) (*Subject, error)
}

// ErrUnsupportedKind informed kind is not supported.
var ErrUnsupportedKind = errors.New("unsupported kind")

// NewProbe instantiates a Probe based on the informed kind, the raw workspaces and params are parsed
// and informed to the new instance.
func NewProbe(cfg *config.Config, kind string, rawWorkspaces, rawParams []string) (Probe, error) {
	c := cfg.GetClientsOrPanic().HTTPClient
	wbs, err := workspaces.Merge([]v1beta1.WorkspaceBinding{}, rawWorkspaces, c)
	if err != nil {
		return nil, err
	}
	printWorkspaceBindings(cfg, wbs)

	prs, err := RawTektonParamsToTypedSlice(rawParams)
	if err != nil {
		return nil, err
	}
	printParams(cfg, prs)

	switch kind {
	case "task":
		return NewTaskProbe(cfg, wbs, prs), nil
	case "pipeline":
		return NewPipelineProbe(cfg, wbs, prs), nil
	}
	return nil, fmt.Errorf("%w: %q", ErrUnsupportedKind, kind)
}
