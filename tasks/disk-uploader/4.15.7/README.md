# KubeVirt Disk Uploader Task

Automates the extraction of disk and uploads it to a container registry, to be used in multiple Kubernetes clusters.

## Prerequisites

VMExport support must be enabled in the feature gates to be available. The [feature gates](https://kubevirt.io/user-guide/cluster_admin/activating_feature_gates/#how-to-activate-a-feature-gate) field in the KubeVirt Custom Resource (CR) must be expanded by adding the VMExport to it.

# Example Scenario

When user runs [KubeVirt Tekton Tasks](https://github.com/kubevirt/kubevirt-tekton-tasks) example pipelines (windows-installer, windows-customize) to prepare Windows disk images - The newly created disk image is only in a single cluster. If user wants to have it in another cluster, then KubeVirt Disk Uploader can be used to push it out of the cluster.

### Parameters

- **EXPORT_SOURCE_KIND**: The name of the export source kind
- **EXPORT_SOURCE_NAME**: The name of the export source
- **VOLUME_NAME**: The volume name (If source kind is PVC, then volume name is equal to source name)
- **IMAGE_DESTINATION**: Destination of the image in container registry
- **PUSH_TIMEOUT**: ContainerDisk push timeout in minutes
- **SECRET_NAME**: Name of the secret which holds credential for container registry

### Usage

Task run using resolver:
```
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
    generateName: disk-uploader-taskrun-resolver-
spec:
    params:
    -   name: EXPORT_SOURCE_KIND
        value: vm
    -   name: EXPORT_SOURCE_NAME
        value: example-vm
    -   name: VOLUME_NAME
        value: example-dv
    -   name: IMAGE_DESTINATION
        value: quay.io/kubevirt/example-vm-exported:latest
    -   name: PUSH_TIMEOUT
        value: 120
    taskRef:
        params:
        -   name: catalog
            value: kubevirt-tekton-tasks
        -   name: type
            value: artifact
        -   name: kind
            value: task
        -   name: name
            value: disk-uploader
        -   name: version
            value: v4.15.7
        resolver: hub

```

### Platforms

The Task can be run on linux/amd64 platform.