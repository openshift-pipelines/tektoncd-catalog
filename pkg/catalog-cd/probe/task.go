package probe

import (
	"context"
	"fmt"

	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/config"

	tkncmdtaskrun "github.com/tektoncd/cli/pkg/cmd/taskrun"
	"github.com/tektoncd/cli/pkg/options"
	tkntaskrun "github.com/tektoncd/cli/pkg/taskrun"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TaskProbe runs and follows a TaskRun instance.
type TaskProbe struct {
	cfg        *config.Config             // global configuration
	workspaces []v1beta1.WorkspaceBinding // task workspaces
	params     []v1beta1.Param            // task params
}

var _ Probe = &TaskProbe{}

// generate generates the TaskRun instance for informed Task name.
func (p *TaskProbe) generate(ctx context.Context, name string) (*v1beta1.TaskRun, error) {
	tr := &v1beta1.TaskRun{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1beta1.SchemeGroupVersion.String(),
			Kind:       "TaskRun",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-run-", name),
		},
		Spec: v1beta1.TaskRunSpec{
			TaskRef: &v1beta1.TaskRef{
				Name: name,
			},
			Workspaces: p.workspaces,
			Params:     p.params,
		},
	}
	if err := tr.Spec.Validate(ctx); err != nil {
		return nil, err
	}
	return tr, nil
}

// Run uses the informed (task) name to generate and create a TaskRun instance, following up its
// progress until completion.
func (p *TaskProbe) Run(ctx context.Context, name string) (*Subject, error) {
	p.cfg.Infof("### Issuing a TaskRun for %q...\n", name)
	tr, err := p.generate(ctx, name)
	if err != nil {
		return nil, err
	}

	cs := p.cfg.GetClientsOrPanic()
	ns := p.cfg.GetNamespace()
	if tr, err = tkntaskrun.Create(cs, tr, metav1.CreateOptions{}, ns); err != nil {
		return nil, err
	}
	p.cfg.Infof("### Issued %q!\n", tr.GetNamespacedName())

	p.cfg.Infof("### Following progress...\n\n")
	opts := &options.LogOptions{
		Follow:      true,
		AllSteps:    true,
		Params:      p.cfg.GetTektonParams(),
		Stream:      p.cfg.Stream,
		Prefixing:   true,
		TaskrunName: tr.GetName(),
	}
	if err = tkncmdtaskrun.Run(opts); err != nil {
		return nil, err
	}

	return NewSubject(
		tr.GetGroupVersionKind().GroupVersion().WithResource("taskruns"),
		tr.GetNamespacedName(),
	), nil
}

// NewTaskProbe instantiate a TaskProbe with requirements.
func NewTaskProbe(
	cfg *config.Config,
	wbs []v1beta1.WorkspaceBinding,
	prs []v1beta1.Param,
) *TaskProbe {
	return &TaskProbe{
		cfg:        cfg,
		workspaces: wbs,
		params:     prs,
	}
}
