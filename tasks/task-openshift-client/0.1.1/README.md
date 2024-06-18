`OpenShift-Cient` Tekton Task
-----------------------

# Abstract

The `openshift-client` Task is the binary for OpenShift CLI that complements kubectl for simplifying deployment and configuration applications on OpenShift.

# Usage

Please, consider the usage example below:

```yaml
---
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata: {}
spec:
  pipelineRef:
    name: task-openshift-client
  params:
    - name: SCRIPT
      value: openShift-CLI-arguments
    - name: VERSION
      value: oc-version

```
You'll need to replace `openShift-CLI-arguments`  with OpenShift CLI arguments that you want to run and replace `oc-version` with the version of OpenShift you want to use.
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
| `manifest_dir` | `true` | The workspace which contains kubernetes manifests which we want to apply on the cluster. |
| `kubeconfig_dir` | `true` | The workspace which contains the the kubeconfig file if in case we want to run the oc command on another cluster. |


## Params

| Param         | Type                       | Default                      | Description                |
| :------------ | :------------------------: | :--------------------------- | :------------------------- |
| `SCRIPT` | `string` | (required) | The OpenShift CLI arguments to run |
| `VERSION` | `string` | (required) | The OpenShift version to use |



[tektonPipelineAuth]: https://tekton.dev/docs/pipelines/auth/#configuring-docker-authentication-for-docker

