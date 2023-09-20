package contract

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ErrTektonResourceUnsupported marks the resource as not supported, as in it's not a
// Kubernetes CRD, or not a Tekton API on supported versions, etc.
var ErrTektonResourceUnsupported = errors.New("tekton resource not supported")

// isResourceSupported inspects the unstructured to assert it's a Tekton resource, and its
// version is supported by this program.
func isResourceSupported(u *unstructured.Unstructured) error {
	group := u.GroupVersionKind().Group
	if group != v1.SchemeGroupVersion.Group {
		return fmt.Errorf("%w: unsupported group %q (expects %q)",
			ErrTektonResourceUnsupported, group, v1.SchemeGroupVersion.Group)
	}
	version := u.GroupVersionKind().Version
	if version != "v1" && version != "v1beta1" {
		return fmt.Errorf("%w: unsupported version %q",
			ErrTektonResourceUnsupported, version)
	}
	kind := u.GetKind()
	if kind != "Task" && kind != "Pipeline" {
		return fmt.Errorf("%w: unsupported kind %q", ErrTektonResourceUnsupported, kind)
	}
	return nil
}

// CalculateSHA256Sum calculates the SHA256 sum of the informed file.
func CalculateSHA256Sum(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
