`Buildah-ns` Tekton Task
-----------------------

# Abstract

The `buildah-ns` Task is meant to build [OCI][OCI] container images without the requirement of container runtime daemon like Docker daemon using [Buildah][Buildah], with enhanced security through user namespace isolation. The Task results contain the image name and the SHA256 image digest.

This task is functionally identical to the `buildah` task but runs with user namespace isolation for enhanced security, providing better isolation between containers at the kernel level.

## Key Differences from `buildah` Task

The `buildah-ns` task differs from the standard `buildah` task in the following ways:

- **Task Name**: `buildah-ns` (vs `buildah`)
- **Annotations**: Includes additional security annotations which enables automatic user namespace mapping:
  - `io.kubernetes.cri-o.userns-mode: "auto"`  
  - `io.openshift.builder: "true"` 
- **Security Model**: Enhanced isolation through user namespace separation

# Usage

Please, consider the usage example below:

```yaml
---
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata: {}
spec:
  pipelineRef:
    name: task-buildah-ns
  params:
    - name: IMAGE
      value: your-image-name
    - name: TLS_VERIFY
      value: true
    - name: VERBOSE
      value: false
  workspaces:
    - name: source
      persistentVolumeClaim:
        claimName: your-pvc-name
```
You'll need to replace `your-image-name`  with the actual name of the image you want to build, and `your-pvc-name`  with the name of the PersistentVolumeClaim where your source code is stored.
In case the Container Registry requires authentication, please consider the [Tekton Pipelines documentation][tektonPipelineAuth]. In a nutshell, you need to create a Kubernetes Secret describing the following attributes:

```bash
kubectl create secret docker-registry imagestreams \
  --docker-server="image-registry.openshift-image-registry.svc:5000" \
  --docker-username="${REGISTRY_USERNAME}" \
  --docker-password="${REGISTRY_TOKEN}"
```

Then make sure the Secret is linked with the Service-Account running the `TaskRun`/`PipelineRun`.

## Workspace

| Name         | Optional                      | Description                      |
| :------------ | :------------------------: | :--------------------------- |
| `source` | `false` | Container build context, like for instnace a application source code followed by a `Containerfile`. |


## Params

| Param         | Type                       | Default                      | Description                |
| :------------ | :------------------------: | :--------------------------- | :------------------------- |
| `IMAGE` | `string` | (required) | Fully qualified source container image name, including tag, to be built by buildah. |
| `CONTAINERFILE_PATH` | `string` | `Containerfile` | Path to the `Containerfile` (or `Dockerfile`) relative to the `source` workspace. |
| `TLS_VERIFY` | `string` | `true` | Sets the TLS verification flags, `true` is recommended. |
| `VERBOSE` | `string` | `false` | Shows a more verbose (debug) output. |
| `SUBDIRECTORY` | `string` | `.` | Relative subdirectory to the `source` Workspace for the build-context. |
| `STORAGE_DRIVER` | `string` | `overlay` | Set buildah storage driver to reflect the currrent cluster node's settings. |
| `BUILD_EXTRA_ARGS` | `string` | `` | Extra parameters passed for the build command when building images. |
| `PUSH_EXTRA_ARGS` | `string` | `` | Extra parameters passed for the push command when pushing images. |
| `SKIP_PUSH` | `string` | `false` | Skip pushing the image to the container registry. |


## Results

| Result        | Description                |
| :------------ | :------------------------- |
| `IMAGE_URL` | Fully qualified image name. |
| `IMAGE_DIGEST` | SHA256 digest of the image just built. |

## Security Features

This task includes enhanced security through user namespace isolation with the following key characteristics:


## UID Mapping Behavior

When the following annotations are applied:

```yaml
io.kubernetes.cri-o.userns-mode: "auto"
io.openshift.builder: "true"
```

the container exhibits the following UID mapping behavior:

- **Inside the container**: The process runs as UID 0 (root within the user namespace)
- **Outside the container**: The process runs as a non-zero UID on the host system

This mapping provides an additional security boundary where:
- The container process appears to have root privileges within its own namespace
- The actual process on the host system runs with reduced privileges
- Potential container escape vulnerabilities are mitigated through namespace isolation

### Security Benefits

- **Kernel-level Isolation**: Provides additional isolation boundaries between containers
- **Enhanced Security**: Helps prevent potential container escape vulnerabilities
- **Privilege Separation**: Maintains the appearance of root access within the container while limiting actual host system privileges

[tektonPipelineAuth]: https://tekton.dev/docs/pipelines/auth/#configuring-docker-authentication-for-docker
[Buildah]: https://github.com/containers/buildah
[OCI]: https://opencontainers.org/ 
