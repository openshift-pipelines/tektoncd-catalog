package linter

import (
	"testing"

	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"

	o "github.com/onsi/gomega"
)

func TestNewLinter(t *testing.T) {
	g := o.NewWithT(t)

	cfg := config.NewConfig()

	l, err := NewLinter(cfg, "../../test/resources/task.yaml")
	g.Expect(err).To(o.Succeed())
	g.Expect(l).NotTo(o.BeNil())

	err = l.Enforce()
	g.Expect(err).To(o.Succeed())
}
