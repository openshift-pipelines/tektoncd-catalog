package assert

import (
	"regexp"
	"testing"

	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

var taskRunResultsTestCases = []struct {
	name    string
	re      regexp.Regexp
	want    string
	wantErr bool
}{{
	name:    "matches first result",
	re:      *regexp.MustCompile(`^FOO=\w+`),
	want:    "FOO=bar",
	wantErr: false,
}, {
	name:    "matches second result",
	re:      *regexp.MustCompile(`^BAR=\w+`),
	want:    "BAR=baz",
	wantErr: false,
}, {
	name:    "doesn't match",
	re:      *regexp.MustCompile(`^404=\w+`),
	want:    "",
	wantErr: true,
}}

func Test_v1beta1TaskRunResultMatchesRegexp(t *testing.T) {
	results := []v1beta1.TaskRunResult{{
		Name:  "FOO",
		Type:  v1beta1.ResultsTypeString,
		Value: *v1beta1.NewArrayOrString("bar"),
	}, {
		Name:  "BAR",
		Type:  v1beta1.ResultsTypeString,
		Value: *v1beta1.NewArrayOrString("baz"),
	}}

	for _, tt := range taskRunResultsTestCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := v1beta1TaskRunResultMatchesRegexp(tt.re, results...)
			if (err != nil) != tt.wantErr {
				t.Errorf("v1beta1TaskRunResultMatchesRegexp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("v1beta1TaskRunResultMatchesRegexp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_taskRunResultMatchesRegexp(t *testing.T) {
	results := []v1.TaskRunResult{{
		Name:  "FOO",
		Type:  v1.ResultsTypeString,
		Value: *v1.NewStructuredValues("bar"),
	}, {
		Name:  "BAR",
		Type:  v1.ResultsTypeString,
		Value: *v1.NewStructuredValues("baz"),
	}}

	for _, tt := range taskRunResultsTestCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := taskRunResultMatchesRegexp(tt.re, results...)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskRunResultMatchesRegexp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("taskRunResultMatchesRegexp() = %v, want %v", got, tt.want)
			}
		})
	}
}
