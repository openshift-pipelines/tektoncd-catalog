package assert

import (
	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/config"

	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"knative.dev/pkg/apis"
)

// printStatusConditions prints out the status conditions.
func printStatusConditions(cfg *config.Config, conditions apis.Conditions) {
	if len(conditions) == 0 {
		return
	}

	cfg.Infof("### Status conditions:\n")
	for _, c := range conditions {
		cfg.Infof("#  - type: %q\n", c.Type)
		cfg.Infof("#    status: %q\n", c.Status)
		cfg.Infof("#    reason: %q\n", c.Reason)
		if c.Severity != "" {
			cfg.Infof("#    severity: %q\n", c.Severity)
		}
		cfg.Infof("#    message: %q\n", c.Message)
	}
}

// printTaskRunResults prints out the TaskRun results.
func printTaskRunResults(cfg *config.Config, results []v1.TaskRunResult) {
	if len(results) == 0 {
		return
	}

	cfg.Infof("### TaskRun results:\n")
	for _, r := range results {
		cfg.Infof("#  - name: %q\n", r.Name)
		cfg.Infof("#    type: %q\n", r.Type)
		cfg.Infof("#    value: %q\n", r.Value.StringVal)
	}
}

// printPipelineRunResults prints out each TaskRun status.
func printPipelineRunResults(
	cfg *config.Config,
	taskRunStatus map[string]*v1beta1.PipelineRunTaskRunStatus,
) {
	if len(taskRunStatus) == 0 {
		return
	}

	cfg.Infof("### PipelineRun results:\n")
	for name, tr := range taskRunStatus {
		cfg.Infof("#  - taskRunName: %q\n", name)
		cfg.Infof("#    results:\n")
		for _, r := range tr.Status.TaskRunResults {
			cfg.Infof("#      - name: %q\n", r.Name)
			cfg.Infof("#        type: %q\n", r.Value.Type)
			cfg.Infof("#        value: %q\n", r.Value.StringVal)
		}
	}
}
