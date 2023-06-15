package probe

import (
	"context"
	"fmt"

	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/config"

	tkncmdpipelinerun "github.com/tektoncd/cli/pkg/cmd/pipelinerun"
	"github.com/tektoncd/cli/pkg/options"
	tknpipelinerun "github.com/tektoncd/cli/pkg/pipelinerun"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PipelineProbe runs a Pipeline until completion.
type PipelineProbe struct {
	cfg        *config.Config             // global configuration
	workspaces []v1beta1.WorkspaceBinding // tekton workspaces
	params     []v1beta1.Param            // tekton params
}

var _ Probe = &PipelineProbe{}

// generate generates a valid PipelineRun instance.
func (p *PipelineProbe) generate(ctx context.Context, name string) (*v1beta1.PipelineRun, error) {
	pr := &v1beta1.PipelineRun{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1beta1.SchemeGroupVersion.String(),
			Kind:       "PipelineRun",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-run-", name),
		},
		Spec: v1beta1.PipelineRunSpec{
			PipelineRef: &v1beta1.PipelineRef{
				Name: name,
			},
			Workspaces: p.workspaces,
			Params:     p.params,
		},
	}
	if err := pr.Spec.Validate(ctx); err != nil {
		return nil, err
	}
	return pr, nil
}

// Run runs the PipelineRun until completion, caputures the coordinates on the returned Subject.
func (p *PipelineProbe) Run(ctx context.Context, name string) (*Subject, error) {
	p.cfg.Infof("### Issuing a PipelineRun for %q...\n", name)
	pr, err := p.generate(ctx, name)
	if err != nil {
		return nil, err
	}

	cs := p.cfg.GetClientsOrPanic()
	ns := p.cfg.GetNamespace()
	if pr, err = tknpipelinerun.Create(cs, pr, metav1.CreateOptions{}, ns); err != nil {
		return nil, err
	}
	p.cfg.Infof("### Issued %q!\n", pr.GetNamespacedName())

	p.cfg.Infof("### Following progress...\n\n")
	opts := &options.LogOptions{
		Follow:          true,
		AllSteps:        true,
		Params:          p.cfg.GetTektonParams(),
		Stream:          p.cfg.Stream,
		Prefixing:       true,
		PipelineName:    name,
		PipelineRunName: pr.GetName(),
	}
	if err = tkncmdpipelinerun.Run(opts); err != nil {
		return nil, err
	}

	return NewSubject(
		pr.GroupVersionKind().GroupVersion().WithResource("pipelineruns"),
		pr.GetNamespacedName(),
	), nil
}

// NewPipelineProbe instantiates the PipelineProbe.
func NewPipelineProbe(
	cfg *config.Config,
	wbs []v1beta1.WorkspaceBinding,
	prs []v1beta1.Param,
) *PipelineProbe {
	return &PipelineProbe{
		cfg:        cfg,
		workspaces: wbs,
		params:     prs,
	}
}
