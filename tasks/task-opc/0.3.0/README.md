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

## Working Example

Here's a complete working example that uses the `opc` task to list pipelines:

**Pipeline:**

```yaml
apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  name: opc-task-pipeline
spec:
  tasks:
    - name: opc-pipeline-list
      taskRef:
        resolver: cluster
        params:
          - name: kind
            value: task
          - name: name
            value: opc
          - name: namespace
            value: openshift-pipelines
      params:
        - name: SCRIPT
          value: "opc $@"
        - name: ARGS
          value:
            - pipeline
            - list
            - -n
            - $(context.pipelineRun.namespace)
```

**PipelineRun:**

```yaml
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: opc-task-run
spec:
  pipelineRef:
    name: opc-task-pipeline
```

**Output:**

```
[opc-pipeline-list : opc] Running Script /scripts/opc-client.sh
[opc-pipeline-list : opc] NAME                AGE              LAST RUN       STARTED         DURATION   STATUS
[opc-pipeline-list : opc] opc-task-pipeline   15 minutes ago   opc-task-run   5 seconds ago   ---        Running
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

