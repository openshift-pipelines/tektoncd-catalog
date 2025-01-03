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

In case multiple images need to be copied, or one image under multiple names/tags, a `url.txt` file can be provided:

```text
docker://docker.io/library/busybox:latest docker://image-registry.openshift-image-registry.svc.cluster.local:5000/my-project/busybox:latest
docker://docker.io/library/busybox:1 docker://image-registry.openshift-image-registry.svc.cluster.local:5000/my-project/busybox:1
```


This file has to be present at the root and mounted under the workspace `images_url`, for example using a ConfigMap:
```yaml
---
apiVersion: tekton.dev/v1beta1
kind: TaskRun
metadata: {}
spec:
  taskRef:
    name: skopeo-copy
  workspaces:
    - name: images_url
      configmap:
        name: configmap-images
```

Referenced config map:

```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: configmap-images
data:
  url.txt: |
    docker://docker.io/library/busybox:latest docker://image-registry.openshift-image-registry.svc.cluster.local:5000/my-project/busybox:latest
    docker://docker.io/library/busybox:1 docker://image-registry.openshift-image-registry.svc.cluster.local:5000/my-project/busybox:1

```

## Workspace

| Name         | Optional | Description                                                                                                                                                                                |
|:-------------|:--------:|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `images_url` |  `true`  | For storing image urls in case of more than one image to copy. It must have a url.txt file at the root path containing a source and a destination image separated by a space on each line. |

## Params

| Param         | Type                       | Default                      | Description                                                                                                                                                          |
| :------------ | :------------------------: | :---------------------------:|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `SOURCE_IMAGE_URL` | `string` | "" | Fully qualified source container image name, including tag, to be copied into `DESTINATION_IMAGE_URL` param. If more than one image needs to be copied, leave empty. |
| `DESTINATION_IMAGE_URL` | `string` | "" | Fully qualified destination container image name, including tag. If more than one image needs to be copied, leave empty.                                                                                                    |
| `SRC_TLS_VERIFY` | `string` | `true` | Sets the TLS verification flags for the source registry, `true` is recommended.                                                                                      |
| `DEST_TLS_VERIFY` | `string` | `true` | Sets the TLS verification flags for the destination registry, `true` is recommended.                                                                                 |
| `VERBOSE` | `string` | `false` | Shows a more verbose (debug) output.                                                                                                                                 |

## Results

| Result        | Description                                                                                                                                                                        |
| :------------ |:-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `SOURCE_DIGEST` | If a single image is copied it contains a source image SHA256 digest. When copying multiple images it contains SHA256 digests of all source images separated by a space.           |
| `DESTINATION_DIGEST` | If a single image is copied it contains a destination image SHA256 digest. When copying multiple images it contains SHA256 digests of all destination images separated by a space. |

[tektonPipelineAuth]: https://tekton.dev/docs/pipelines/auth/#configuring-docker-authentication-for-docker
[containersSkopeo]: https://github.com/containers/skopeo
