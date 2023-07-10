package render

import (
	"testing"

	o "github.com/onsi/gomega"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
)

func TestNewMarkdow(t *testing.T) {
	g := o.NewWithT(t)

	cfg := config.NewConfig()

	m, err := NewMarkdown(cfg, "../../../test/resources/task.yaml")
	g.Expect(err).To(o.Succeed())
	g.Expect(m).NotTo(o.BeNil())

	err = m.Render()
	g.Expect(err).To(o.Succeed())
}
