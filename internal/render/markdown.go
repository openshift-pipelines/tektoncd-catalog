package render

import (
	"text/template"

	"github.com/openshift-pipelines/tektoncd-catalog/internal/config"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/linter"
	"github.com/openshift-pipelines/tektoncd-catalog/internal/resource"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	_ "embed"
)

// Markdown renders a Tekton resource workspaces, params and results as markdown tables.
type Markdown struct {
	cfg *config.Config             // global configuration
	u   *unstructured.Unstructured // object instance
}

//go:embed tekton.md.tpl
var markdownTemplate []byte

// templateInputs extracts the inputs for the template.
func (m *Markdown) templateInputs() (map[string][]interface{}, error) {
	inputs := map[string][]interface{}{}
	var err error

	for _, attribute := range []string{"workspaces", "params", "results"} {
		inputs[attribute], err = linter.GetNestedSlice(m.u, "spec", attribute)
		if err != nil {
			return nil, err
		}
	}
	return inputs, nil
}

// Render instantiate a new template using the local functions to render the resource as markdown.
func (m *Markdown) Render() error {
	tpl, err := template.New("markdown").Funcs(templateFuncMap).Parse(string(markdownTemplate))
	if err != nil {
		return err
	}
	inputs, err := m.templateInputs()
	if err != nil {
		return err
	}
	return tpl.Execute(m.cfg.Stream.Out, inputs)
}

// NewMarkdown instantiates the markdown render by decoding the informed resource file.
func NewMarkdown(cfg *config.Config, resourceFile string) (*Markdown, error) {
	u, err := resource.ReadAndDecodeResourceFile(resourceFile)
	if err != nil {
		return nil, err
	}
	return &Markdown{cfg: cfg, u: u}, nil
}
