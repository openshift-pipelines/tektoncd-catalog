package probe

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// Subject defines the coordinates of the Kubernetes (Tekton) resource to be "subject" of probing and
// assertion.
type Subject struct {
	GVR      schema.GroupVersionResource // kubernetes GVR
	Fullname types.NamespacedName        // resource namespace and name
}

// Name short access to subject's namespace.
func (s *Subject) Namespace() string {
	return s.Fullname.Namespace
}

// Name short access to subject's name.
func (s *Subject) Name() string {
	return s.Fullname.Name
}

// NewSubject instantiate the Subject.
func NewSubject(gv schema.GroupVersionResource, fullname types.NamespacedName) *Subject {
	return &Subject{
		GVR:      gv,
		Fullname: fullname,
	}
}
