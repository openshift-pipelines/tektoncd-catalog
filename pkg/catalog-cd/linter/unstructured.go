package linter

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func GetNestedSlice(u *unstructured.Unstructured, fields ...string) ([]interface{}, error) {
	slice, found, err := unstructured.NestedSlice(u.Object, fields...)
	if err != nil {
		return nil, err
	}
	if !found {
		return []interface{}{}, nil
	}
	return slice, nil
}
