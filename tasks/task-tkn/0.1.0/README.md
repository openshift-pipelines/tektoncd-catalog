`Tkn` Tekton Task
-----------------------

# Abstract

The `tkn` Task is a binary to perform operations on Tekton resources using tkn.

# Usage

Please, consider the usage example below:

```yaml
---
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata: {}
spec:
  pipelineRef:
    name: task-tkn
  params:
    - name: SCRIPT
      value: "tkn $@"
    - name: ARGS
      value: tkn-command-you-want-to-execute

```
You'll need to replace `tkn-command-you-want-to-execute`  with tkn CLI arguments based on what operation you want to perform on the Tekton resources.
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
| `kubeconfig_dir` | `true` | An optional workspace that allows you to provide a .kube/config file for tkn to access the cluster. |


## Params

| Param         | Type                       | Default                      | Description                |
| :------------ | :------------------------: | :--------------------------- | :------------------------- |
| `SCRIPT` | `string` | (required) | tkn CLI script to execute |
| `ARGS` | `array` | (required) | tkn CLI arguments to run |



[tektonPipelineAuth]: https://tekton.dev/docs/pipelines/auth/#configuring-docker-authentication-for-docker

