Containers Tekton Tasks
-----------------------

# Abstract

Describes the Tekton Tasks supporting Skopeo-Copy

# `skopeo-copy` Tekton Task

The `skopeo-copy` Task is meant to replicate a container image from the `SOURCE_IMAGE_URL` registry to the `DESTINATION_IMAGE_URL` using [Skopeo][containersSkopeo], the Task results contain the SHA256 digests.

# Usage

Please, consider the usage example below:

```yaml
---
apiVersion: tekton.dev/v1beta1
kind: TaskRun
metadata: {}
spec:
  taskRef:
    name: skopeo-copy
  params:
    - name: SOURCE_IMAGE_URL
      value: docker://docker.io/busybox:latest
    - name: DESTINATION_IMAGE_URL
      value: docker://image-registry.openshift-image-registry.svc:5000/task-containers/busybox:latest
```

In case the Container Registry requires authentication, please consider the [Tekton Pipelines documentation][tektonPipelineAuth]. In a nutshell, you need to create a Kubernetes Secret describing the following attributes:

```bash
kubectl create secret docker-registry imagestreams \
  --docker-server="image-registry.openshift-image-registry.svc:5000" \
  --docker-username="${REGISTRY_USERNAME}" \
  --docker-password="${REGISTRY_TOKEN}"
```

Then make sure the Secret is linked with the Service-Account running the `TaskRun`/`PipelineRun`.

## Params

| Param         | Type                       | Default                      | Description                |
| :------------ | :------------------------: | :--------------------------- | :------------------------- |
| `SOURCE_IMAGE_URL` | `string` | (required) | Fully qualified source container image name, including tag, to be copied into `DESTINATION_IMAGE_URL` param. |
| `DESTINATION_IMAGE_URL` | `string` | (required) | Fully qualified destination container image name, including tag. |
| `SRC_TLS_VERIFY` | `string` | `true` | Sets the TLS verification flags for the source registry, `true` is recommended. |
| `DEST_TLS_VERIFY` | `string` | `true` | Sets the TLS verification flags for the destination registry, `true` is recommended. |
| `VERBOSE` | `string` | `false` | Shows a more verbose (debug) output. |

## Results

| Result        | Description                |
| :------------ | :------------------------- |
| `SOURCE_DIGEST` | Source image SHA256 digest. |
| `DESTINATION_DIGEST` | Destination image SHA256 digest. |

[tektonPipelineAuth]: https://tekton.dev/docs/pipelines/auth/#configuring-docker-authentication-for-docker
[containersSkopeo]: https://github.com/containers/skopeo
