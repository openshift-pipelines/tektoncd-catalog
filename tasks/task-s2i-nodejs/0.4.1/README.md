Source-to-Image Tekton Tasks (`s2i`)
------------------------------------

# Abstract

Describes the Tekton Tasks supporting Source-to-Image for various ecosystems.

# `s2i` Builder Images

This section we explain each language ecosystem comes with a predefined set of builder images, supported by Red Hat.

The `s2i` Task helps in building reproducible container images from source code i.e. task for supporting s2i functionality.

The s2i Task has been customized with builder images specific to various languages and have been named appropriately as follows:

| Task Name  | Builder Image Used                                                     |
| ---------- | ---------------------------------------------------------------------- |
| s2i-python | http://registry.access.redhat.com/ubi8/python-39:latest                |
| s2i-go     | http://registry.access.redhat.com/ubi8/go-toolset:1.19.10-3            |
| s2i-java   | http://registry.access.redhat.com/ubi8/openjdk-11:latest               |
| s2i-dotnet | http://registry.access.redhat.com/ubi8/dotnet-60:6.0-37.20230802191230 |
| s2i-php    | http://registry.access.redhat.com/ubi9/php-81:1-29                     |
| s2i-nodejs | http://registry.access.redhat.com/ubi8/nodejs-18:latest                |
| s2i-perl   | http://registry.access.redhat.com/ubi9/perl-532:1-91                   |
| s2i-ruby   | http://registry.access.redhat.com/ubi9/ruby-31:1-50                    |

In case, the above builder images associated with the languages aren’t satisfactory for your source code, you can change it using appropriate parameter.

# Usage

Please, consider the usage example below:

```yaml
---
apiVersion: tekton.dev/v1beta1
kind: TaskRun
metadata: {}
spec:
  taskRef:
    name: s2i-python
  params:
    - name: IMAGE
      value: registry.registry.svc.cluster.local:32222/task-containers/task-s2i-python:latest
```

In case the Container Registry requires authentication, please consider the [Tekton Pipelines documentation][tektonPipelineAuth]. In a nutshell, you need to create a Kubernetes Secret describing the following attributes:

```bash
kubectl create secret docker-registry imagestreams \
  --docker-server="image-registry.openshift-image-registry.svc:5000" \
  --docker-username="${REGISTRY_USERNAME}" \
  --docker-password="${REGISTRY_TOKEN}"
```

Then make sure the Secret is linked with the Service-Account running the `TaskRun`/`PipelineRun`.

## Workspaces

All of the s2i tasks use the `source` workspace which is meant to contain the Application source code, which acts as the build context for S2I workflow.


## Params

| Param             | Type   | Default                  | Description                                                               |
| ----------------- | ------ | ------------------------ | ------------------------------------------------------------------------- |
| IMAGE             | string | (required)               | Fully qualified container image name to be built by s2i                   |
| IMAGE_SCRIPTS_URL | string | image:///usr/libexec/s2i | URL containing the default assemble and run scripts for the builder image |
| ENV_VARS          | array  | []                       | Array containing string of Environment Variables as "KEY=VALUE”           |
| SUBDIRECTORY      | string | .                        | Relative subdirectory to the source Workspace for the build-context.      |
| STORAGE_DRIVER    | string | overlay                  | Set buildah storage driver to reflect the currrent cluster node's         |
| settings.         |
| BUILD_EXTRA_ARGS  | string |                          | Extra parameters passed for the build command when building images.       |
| PUSH_EXTRA_ARGS   | string |                          | Extra parameters passed for the push command when pushing images.         |
| SKIP_PUSH         | string | false                    | Skip pushing the image to the container registry.                         |
| TLS_VERIFY        | string | true                     | Sets the TLS verification flag, true is recommended.                      |
| VERBOSE           | string | false                    | Turns on verbose logging, all commands executed will be printed out.      |

## Results

| Result       | Description                     |
| ------------ | ------------------------------- |
| IMAGE_URL    | Fully qualified image name.     |
| IMAGE_DIGEST | Digest of the image just built. |
