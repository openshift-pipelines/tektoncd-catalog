package linter

import (
	"os"

	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
)

// ReadAndDecodeResourceFile reads the informed file and decode contents using Tekton's Kubernetes
// schema, returning a Unstructured instance.
func ReadAndDecodeResourceFile(resource string) (*unstructured.Unstructured, error) {
	payload, err := os.ReadFile(resource)
	if err != nil {
		return nil, err
	}

	runtimeScheme := runtime.NewScheme()

	if err = scheme.AddToScheme(runtimeScheme); err != nil {
		return nil, err
	}
	if err = v1beta1.AddToScheme(runtimeScheme); err != nil {
		return nil, err
	}
	if err = v1.AddToScheme(runtimeScheme); err != nil {
		return nil, err
	}

	obj, _, err := serializer.NewCodecFactory(runtimeScheme).
		UniversalDeserializer().
		Decode(payload, nil, nil)
	if err != nil {
		return nil, err
	}

	u := unstructured.Unstructured{Object: map[string]interface{}{}}
	if u.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(obj); err != nil {
		return nil, err
	}
	return &u, nil
}
