## `Kn-Apply` Tekton Task

# Abstract

The `kn-apply` deploys a given image to a Knative Service using kn command line interface.

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
    name: kn-apply
  params:
    - name: SERVICE
      value: service-name
    - name: IMAGE
      value: image-name
```

You'll need to replace `service-name` & `image-name` with values of the Knative Service & Image to be deployed respectively.

In case the Container Registry requires authentication, please consider the [Tekton Pipelines documentation][tektonPipelineAuth]. In a nutshell, you need to create a Kubernetes Secret describing the following attributes:

```bash
kubectl create secret docker-registry imagestreams \
  --docker-server="image-registry.openshift-image-registry.svc:5000" \
  --docker-username="${REGISTRY_USERNAME}" \
  --docker-password="${REGISTRY_TOKEN}"
```

Then make sure the Secret is linked with the Service-Account running the `TaskRun`/`PipelineRun`.

## Params

| Param     |   Type   | Default    | Description                                                  |
| :-------- | :------: | :--------- | :----------------------------------------------------------- |
| `SERVICE` | `string` | (required) | Knative Service name to which the given image is deployed to |
| `IMAGE`   | `string` | (required) | The image that needs to be deployed to the Knative Service   |

[tektonPipelineAuth]: https://tekton.dev/docs/pipelines/auth/#configuring-docker-authentication-for-docker
