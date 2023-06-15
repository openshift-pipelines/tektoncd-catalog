package assert

import (
	"errors"
	"fmt"
	"regexp"

	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
)

var (
	// ErrUnmatchedRegexp regular expression doesn't match any results.
	ErrUnmatchedRegexp = errors.New("unmatched regexp")

	// ErrStepFailed indicates one or more steps (status condition) failed.
	ErrStepFailed = errors.New("failed step")
)

// v1beta1TaskRunResultMatchesRegexp asserts one of the informed results matches the regexp.
func v1beta1TaskRunResultMatchesRegexp(
	re regexp.Regexp,
	results ...v1beta1.TaskRunResult,
) (string, error) {
	for _, result := range results {
		kv := fmt.Sprintf("%s=%s", result.Name, result.Value.StringVal)
		if re.MatchString(kv) {
			return kv, nil
		}
	}
	return "", fmt.Errorf("%w: '%v'", ErrUnmatchedRegexp, re.String())
}

// taskRunResultMatchesRegexp asserts one of the informed results matches the regexp.
func taskRunResultMatchesRegexp(re regexp.Regexp, results ...v1.TaskRunResult) (string, error) {
	for _, result := range results {
		kv := fmt.Sprintf("%s=%s", result.Name, result.Value.StringVal)
		if re.MatchString(kv) {
			return kv, nil
		}
	}
	return "", fmt.Errorf("%w: '%v'", ErrUnmatchedRegexp, re.String())
}

// allStepsSucceeded asserts all the status conditions have succeeded.
func allStepsSucceeded(conditions apis.Conditions) error {
	for _, condition := range conditions {
		if condition.Type != apis.ConditionSucceeded || condition.Status != corev1.ConditionTrue {
			continue
		}
		return nil
	}
	return fmt.Errorf("%w: %d conditions inspected", ErrStepFailed, len(conditions))
}
