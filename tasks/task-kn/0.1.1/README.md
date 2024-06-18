## `Kn` Tekton Task

# Abstract

The `kn` Task performs operations on Knative resources (services, revisions, routes) using kn CLI.

# Usage

Please, consider the usage example below:

```yaml
---
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
  name: example-taskrun
spec:
  taskRef:
    name: kn
  params:
    - name: ARGS
      value:
        - help
```

Here `ARGS` param accepts an array of arguments for usage with kn command.

In case the Container Registry requires authentication, please consider the [Tekton Pipelines documentation][tektonPipelineAuth]. In a nutshell, you need to create a Kubernetes Secret describing the following attributes:

```bash
kubectl create secret docker-registry imagestreams \
  --docker-server="image-registry.openshift-image-registry.svc:5000" \
  --docker-username="${REGISTRY_USERNAME}" \
  --docker-password="${REGISTRY_TOKEN}"
```

Then make sure the Secret is linked with the Service-Account running the `TaskRun`/`PipelineRun`.

## Params

| Param  |  Type   | Default    | Description             |
| :----- | :-----: | :--------- | :---------------------- |
| `ARGS` | `array` | (required) | kn CLI arguments to run |

[tektonPipelineAuth]: https://tekton.dev/docs/pipelines/auth/#configuring-docker-authentication-for-docker
