package assert

import (
	"context"
	"regexp"

	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/probe"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/tektoncd/pipeline/pkg/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PipelineRunAssert asserts the elements of a PipelineRun instance.
type PipelineRunAssert struct {
	cfg           *config.Config                               // global configuration
	pr            *v1beta1.PipelineRun                         // pipelinerun instance (subject)
	taskRunStatus map[string]*v1beta1.PipelineRunTaskRunStatus // taskrun children instances
}

var _ Assert = &PipelineRunAssert{}

// Status asserts the instance status conditions.
func (a *PipelineRunAssert) Status() error {
	return allStepsSucceeded(a.pr.Status.GetConditions())
}

// Results asserts the PipelineRun children TaskRun instances.
func (a *PipelineRunAssert) Results(rules []regexp.Regexp) error {
	results := []v1beta1.TaskRunResult{}
	for _, tr := range a.taskRunStatus {
		results = append(results, tr.Status.TaskRunResults...)
	}

	a.cfg.Infof("### Asserting results:\n")
	for _, re := range rules {
		a.cfg.Infof("#  - regexp: '%v'\n", re.String())
		result, err := v1beta1TaskRunResultMatchesRegexp(re, results...)
		if err != nil {
			a.cfg.Infof("#    match: \"\"\n")
			return err
		}
		a.cfg.Infof("#    match: %q\n", result)
	}
	return nil
}

// NewPipelineRunAssert instantiates the PipelineRunAssert with a up-to-date PipelineRun instance and
// the children TaskRun statuses.
func NewPipelineRunAssert(
	ctx context.Context,
	cfg *config.Config,
	subject *probe.Subject,
) (*PipelineRunAssert, error) {
	cs := cfg.GetClientsOrPanic().Tekton
	ns := subject.Namespace()

	// retrieving the latest changes for the informed subject
	pr, err := cs.TektonV1beta1().PipelineRuns(ns).Get(ctx, subject.Name(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	printStatusConditions(cfg, pr.Status.GetConditions())

	// going after the status of the TaskRun children instances
	taskRunStatus, _, err := status.GetFullPipelineTaskStatuses(ctx, cs, ns, pr)
	if err != nil {
		return nil, err
	}
	printPipelineRunResults(cfg, taskRunStatus)

	return &PipelineRunAssert{
		cfg:           cfg,
		pr:            pr,
		taskRunStatus: taskRunStatus,
	}, nil
}
