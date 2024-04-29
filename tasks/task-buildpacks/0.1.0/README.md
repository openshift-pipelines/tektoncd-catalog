<p align="center">
    <a alt="Test Workflow" href="https://github.com/openshift-pipelines/task-buildpacks/actions/workflows/test.yaml">
        <img alt="GitHub Test Workflow Status" src="https://img.shields.io/github/actions/workflow/status/openshift-pipelines/task-buildpacks/test.yaml?label=test">
    </a>
    <a alt="Latest Release" href="https://github.com/openshift-pipelines/task-buildpacks/releases/latest">
        <img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/openshift-pipelines/task-buildpacks">
    </a>
</p>

`buildpacks` Tekton Task
------------------------

Tekton Task for the [Buildpacks][buildpacksIO], to build a application from source code into a container image.

# Usage

Please consider a example usage to build [Paketo's Node.js sample][paketoNodejsSample] below. The `Pipeline` uses [task-git (0.0.1)][taskGit] to clone the repository:

```yaml
---
apiVersion: tekton.dev/v1beta1
kind: Pipeline
metadata: {}
spec:
  workspaces:
    - name: source

  params:
    - name: IMAGE
      type: string
      description: Fully qualified image name, including tag

  tasks:
    # Cloning the git repository (URL) onto the "source" workspace.
    - name: git
      taskRef:
        name: git
      workspaces:
        - name: output
          workspace: source
      params:
        - name: URL
          value: https://github.com/paketo-buildpacks/samples.git

    # Running the buildpacks' CNB on the source code to build the container
    # image (IMAGE param).
    - name: buildpacks
      taskRef:
        name: buildpacks
      runAfter:
        - git
      workspaces:
        - name: source
          workspace: source
      params:
        - name: SUBDIRECTORY
          value: nodejs/npm
        - name: IMAGE
          value: $(params.IMAGE)
```

Please consider the [Workspaces](#workspaces), [Parameters](#parameters) and [Results](#results) described below.

## Container Registry Authentication

The `buildpacks` Task needs credentials to communicate with the Container Registries to publish the images created by the CNB (`IMAGE` param), and also, to optionally store cached data (`CACHE_IMAGE` param).

Please consider Tekton Pipeline documentation to setup the [Container Registry authentication][tektonContainerRegistryAuth].

# Workspaces

| Workspace      | Optional                           | Description                |
| :------------- | :--------------------------------: | :------------------------- |
| `source`  | `false` | Application source-code. |
| `cache`  | `true` | Cache directory, alternative to the `CACHE_IMAGE` param. |
| `bindings`  | `true` | Extra bindings, CA certificate bundle files. |

## `source`

Contains the application source code, please consider [task-git][taskGit] to handle the Git repository clone. The param `SUBDIRECTORY` defines a relative path on this Workspace to be used as application's source code alternatively, as the [usage](#usage) example shows.

## `cache`

Contains the build cache to be reused on subsequent executions. Caching can also be achieved using the `CACHE_IMAGE` param to store the cache data as a container image.

Considering the nature of the cached data, it's recommended to use a [Kubernetes Persistent Volume Claim][tektonWorkspaceVolume].

## `bindings`

Workspace for the extra binding certificates needed for the application runtime. The files present on this Workspace matching `BINDINGS_GLOB` will be copied into `SERVICE_BINDING_ROOT`, where the CNB picks up the binding files.

For instance, the following Secret carries a `example.pem` file:

```sh
kubectl create secret generic buildpacks-ex --from-file="example.pem=/path/to/example.pem"
```

Given the param `BINDINGS_GLOB` is set to pick up `*.pem` files, then you can [mount the "buildpacks-ex"][tektonWorkspaceSecret] Secret on the Task's `bindings` workspace.

# Params

| Param         | Type                       | Default                      | Description                |
| :------------ | :------------------------: | :--------------------------- | :------------------------- |
| `IMAGE` | `string` | (required) | Application's container image name, and tag. |
| `BUILDER_IMAGE` | `string` | `docker.io/paketobuildpacks/builder:base` | Cloud Native Builder (CNB) container image name (and tag). |
| `CNB_PLATFORM_API` | `string` | `0.11` | Lifecycle platform API compatibility version. |
| `SUBDIRECTORY` | `string` | "" (empty) | Alternative `CNB_APP_DIR` directory, relative to the "source" Workspace. |
| `ENV_VARS` | `array` | `[]` (empty) | Environment variables to set during "build-time". |
| `PROCESS_TYPE` | `string` | `web` | Application process type. |
| `BINDINGS_GLOB` | `string` | `*.pem` | Extra binding file name(s) (glob expression) present on the `bindings` Workspace to be copied into `SERVICE_BINDING_ROOT` directory. |
| `RUN_IMAGE` | `string` | "" (empty) | Reference to a run image to use. |
| `CACHE_IMAGE` | `string` | "" (empty) | The name of the persistent cache image (when  cache workspace is not provided). |
| `SKIP_RESTORE` | `string` | `false` | Do not write layer metadata or restore cached layers. |
| `USER_ID` | `string` | `1000` | CNB container image user-id (UID). |
| `GROUP_ID` | `string` | `1000` | CNB container image group-id (GID). |
| `VERBOSE` | `string` | `false` | Turns on verbose logging, all commands executed will be printed out. |

# Results

| Result        | Description                |
| :------------ | :------------------------- |
| `IMAGE_DIGEST` | Reported `IMAGE` digest. |
| `IMAGE_URL` | Reported fully qualified container image name. |

[buildpacksIO]: https://buildpacks.io/
[paketoNodejsSample]: https://github.com/paketo-buildpacks/samples/tree/main/nodejs/npm
[taskGit]: https://github.com/openshift-pipelines/task-git
[tektonContainerRegistryAuth]: https://github.com/openshift-pipelines/setup-tektoncd#container-registry
[tektonWorkspaceSecret]: https://tekton.dev/docs/pipelines/workspaces/#secret
[tektonWorkspaceVolume]: https://tekton.dev/docs/pipelines/workspaces/#specifying-volumesources-in-workspaces
