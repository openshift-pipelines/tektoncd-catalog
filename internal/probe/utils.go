package probe

import (
	"github.com/tektoncd/cli/pkg/params"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// RawTektonParamsToTypedSlice parses the informed slice of strings as typed Tekton Params.
func RawTektonParamsToTypedSlice(rawParams []string) ([]v1beta1.Param, error) {
	kv, err := params.ParseParams(rawParams)
	if err != nil {
		return nil, err
	}

	paramSlice := []v1beta1.Param{}
	for k, v := range kv {
		paramSlice = append(paramSlice, v1beta1.Param{
			Name: k,
			Value: v1beta1.ParamValue{
				Type:      v1beta1.ParamTypeString,
				StringVal: v,
			},
		})
	}
	return paramSlice, nil
}
