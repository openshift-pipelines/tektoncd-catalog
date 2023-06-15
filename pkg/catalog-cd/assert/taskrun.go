package assert

import (
	"regexp"

	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/config"
	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/probe"

	tkntaskrun "github.com/tektoncd/cli/pkg/taskrun"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

// TaskRunAssert asserts a TaskRun instance.
type TaskRunAssert struct {
	cfg *config.Config // global configuration
	tr  *v1.TaskRun    // taskrun instance (subject)
}

var _ Assert = &TaskRunAssert{}

// Status assert the TaskRun status conditions.
func (a *TaskRunAssert) Status() error {
	return allStepsSucceeded(a.tr.Status.GetConditions())
}

// Results asserts the TaskRun results against the informed regular expressions.
func (a *TaskRunAssert) Results(rules []regexp.Regexp) error {
	a.cfg.Infof("### Asserting results:\n")
	for _, re := range rules {
		a.cfg.Infof("#  - regexp: '%v'\n", re.String())
		result, err := taskRunResultMatchesRegexp(re, a.tr.Status.Results...)
		if err != nil {
			a.cfg.Infof("#    match: \"\"\n")
			return err
		}
		a.cfg.Infof("#    match: %q\n", result)
	}
	return nil
}

// NewTaskRunAssert instantiate the TaskRunAssert by loading the TaskRun (subject) resource.
func NewTaskRunAssert(cfg *config.Config, subject *probe.Subject) (*TaskRunAssert, error) {
	cs := cfg.GetClientsOrPanic()
	tr, err := tkntaskrun.GetTaskRun(subject.GVR, cs, subject.Name(), subject.Namespace())
	if err != nil {
		return nil, err
	}
	printStatusConditions(cfg, tr.Status.GetConditions())
	printTaskRunResults(cfg, tr.Status.Results)

	return &TaskRunAssert{
		cfg: cfg,
		tr:  tr,
	}, nil
}
