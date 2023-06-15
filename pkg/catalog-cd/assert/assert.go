package assert

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/config"
	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/probe"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// Assert defines the Assert role, responsible for asserting the final status of the different type
// of resources asserted by this application.
type Assert interface {
	// Status asserts the instance status conditions.
	Status() error

	// Results asserts the instance results against the slice of regular expressions.
	Results(_ []regexp.Regexp) error
}

// ErrUnsupportedGV the resource group version informed (subject) is not supported.
var ErrUnsupportedGV = errors.New("unsupported resource group-version")

// NewAssert proxy the Assert instantiation based on the subject resource group-version (GVR).
func NewAssert(ctx context.Context, cfg *config.Config, subject *probe.Subject) (Assert, error) {
	switch subject.GVR {
	case v1beta1.SchemeGroupVersion.WithResource("taskruns"):
		return NewTaskRunAssert(cfg, subject)
	case v1beta1.SchemeGroupVersion.WithResource("pipelineruns"):
		return NewPipelineRunAssert(ctx, cfg, subject)
	default:
		return nil, fmt.Errorf("%w: %q (%q)", ErrUnsupportedGV, subject.Fullname, subject.GVR)
	}
}
