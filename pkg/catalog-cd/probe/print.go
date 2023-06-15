package probe

import (
	"github.com/openshift-pipelines/tektoncd-catalog/pkg/catalog-cd/config"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// printWorkspaceBindings prints out the informed workspace bindings slice.
func printWorkspaceBindings(cfg *config.Config, workspaces []v1beta1.WorkspaceBinding) {
	if len(workspaces) == 0 {
		return
	}

	cfg.Infof("### Workspaces:\n")
	for _, w := range workspaces {
		cfg.Infof("#  - name: %q\n", w.Name)
		if w.EmptyDir != nil {
			cfg.Infof("#    emptyDir: %q\n", w.EmptyDir.String())
		}
		if w.ConfigMap != nil {
			cfg.Infof("#    configmap: %q\n", w.ConfigMap.String())
		}
		if w.Secret != nil {
			cfg.Infof("#    secret: %q\n", w.Secret.String())
		}
		if w.Projected != nil {
			cfg.Infof("#    projected: %q\n", w.Projected.String())
		}
		if w.PersistentVolumeClaim != nil {
			cfg.Infof("#    persistentVolumeClaim: %q\n", w.PersistentVolumeClaim.String())
		}
		if w.VolumeClaimTemplate != nil {
			cfg.Infof("#    volumeClaimTemplate: %q\n", w.VolumeClaimTemplate.String())
		}
		if w.SubPath != "" {
			cfg.Infof("#    subpath: %q\n", w.SubPath)
		}
	}
}

// printParams prints out the informed Params slice.
func printParams(cfg *config.Config, params []v1beta1.Param) {
	if len(params) == 0 {
		return
	}

	cfg.Infof("### Params:\n")
	for _, p := range params {
		cfg.Infof("#  - name: %q\n", p.Name)
		cfg.Infof("#    type: %q\n", p.Value.Type)
		cfg.Infof("#    value: %q\n", p.Value.StringVal)
	}
}
