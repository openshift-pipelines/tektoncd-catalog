`opc` Tekton Task
-----------------------

# Abstract

The `opc` Task makes it easy to work with Tekton resources in OpenShift Pipelines.

# Usage

Please, consider the usage example below:

```yaml
---
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata: {}
spec:
  pipelineRef:
    name: task-opc
  params:
    - name: SCRIPT
      value: "opc $@"
    - name: ARGS
      value: opc-command-you-want-to-execute

```
You'll need to replace `opc-command-you-want-to-execute`  with opc CLI arguments based on what operation you want to perform with the Tekton resources in OpenShift Pipelines.
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
| `kubeconfig_dir` | `true` | An optional workspace that allows you to provide a .kube/config file for opc to access the cluster. |


## Params

| Param         | Type                       | Default                      | Description                |
| :------------ | :------------------------: | :--------------------------- | :------------------------- |
| `SCRIPT` | `string` | (required) | opc CLI script to execute |
| `ARGS` | `array` | (required) | opc CLI arguments to run |



[tektonPipelineAuth]: https://tekton.dev/docs/pipelines/auth/#configuring-docker-authentication-for-docker

