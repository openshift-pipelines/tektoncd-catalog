# Red Hat Tekton Catalog
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/redhat-tekton-tasks)](https://artifacthub.io/packages/search?repo=redhat-tekton-tasks) [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/redhat-tekton-pipelines)](https://artifacthub.io/packages/search?repo=redhat-tekton-pipelines)

# Introduction

This is the home of the catalog of Red Hat Tekton resources. The repository contains a catalog of `Task` resources (someday `Pipeline`s and more), designed to be reusable in many Pipelines, [authored and supported by Red Hat](./docs/ecosystem-team.md).

These `Task` and `Pipeline` are coming from the external repositories releases, maintained by different teams from Red Hat and partners. See [here](https://github.com/openshift-pipelines/tektoncd-catalog/blob/main/externals.yaml) to know where they are pulled from.

As of today, they are indexed by [ArtifactHub][artifactHub] in several catalogs:
- [redhat-tekton-tasks][artifactHubRedHatTasks] [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/redhat-tekton-tasks)](https://artifacthub.io/packages/search?repo=redhat-tekton-tasks)
- [redhat-tekton-pipelines][artifactHubRedHatPipelines] [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/redhat-tekton-pipelines)](https://artifacthub.io/packages/search?repo=redhat-tekton-pipelines)
- [redhat-tekton-experimental-tasks][artifactHubRedHatExperimentalTasks] [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/redhat-tekton-experimental-tasks)](https://artifacthub.io/packages/search?repo=redhat-tekton-experimental-tasks)
- [redhat-tekton-experimental-pipelines][artifactHubRedHatExperimentalPipelines] [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/redhat-tekton-experimental-pipelines)](https://artifacthub.io/packages/search?repo=redhat-tekton-experimental-pipelines)

The `main` branch of the repository contains the configuration and tooling to maintain the catalog. The [`p` branch][pBranch] is where the the catalog gets "persisted" and should be consumed.

Each `Task` is provided in a separate directory along with `README.md` and Kubernetes manifests, you can choose which `Task`s to install on your cluster. A directory can hold one task and multiple versions.

The layout of this repository ([branch `p`][pBranch]) the follows the [Tekton Catalog Organization (TEP-0003)][TEP0003].

# Usage

This section explains how to use the catalog ([`p` branch][pBranch]) with the help of various tools like [Tekton Resolvers][tektonResolvers] as well as [Pipelines as Code][pipelineAsCode]. 

## Tekton Git Resolver

[Tekton Git Resolver][tektonGitResolver] retrives the resources directly from this repository [`p` branch][pBranch], like the following `TaskRun` example:

```yml
---
apiVersion: tekton.dev/v1beta1
kind: TaskRun
metadata:
  name: task-git
spec:
  taskRef:
    resolver: git
    params:
      - name: url
        value: https://github.com/openshift-pipelines/tektoncd-catalog
      - name: revision
        value: p
      - name: pathInRepo
        value: tasks/task-git/0.1.0/task-git.yaml
```

The same approach work for [`PipelineRun` resources][tektonGitResolverPipeline].

## Pipeline As Code

Make sure [Pipeline-as-Code (PaC)][pipelineAsCode] is installed and ready on your cluster. For [OpenShift Pipelines][openshiftPipelines], you can define a [`TektonConfig`][openshiftPipelinesConfig] with this catalog by default, i.e:

```yml
---
apiVersion: operator.tekton.dev/v1alpha1
kind: TektonConfig
metadata:
  name: config
spec:
  pipeline:
    git-resolver-config:
      default-url: https://github.com/openshift-pipelines/tektoncd-catalog
      default-revision: p
      fetch-timeout: 1m 
```

Then, on the [repositories being watched by PaC][pipelineAsCodeRepository] you can consume this catalog resources like the following example:

```yml
---
apiVersion: tekton.dev/v1beta1
kind: TaskRun
metadata:
  name: task-git
spec:
  taskRef:
    resolver: git
    params:
      - name: pathInRepo
        value: tasks/task-git/0.1.0/task-git.yaml
```

Skipping the repository `url` and `revision` from `.spec.taskRef.params[]`.

Alternatively the same notation supported described on [Tekton Git Resolver](#tekton-git-resolver) section is supported by PaC.

# Contributing

The Tekton resources on this repository are following the [policies defined here](./docs/lint.md), on a adding new [external repositories references](./externals.yaml) please observe these linting rules first.

External repositories must define a [`catalog.yaml` manifest](./docs/catalog.md), which describes all the Tekton resource on the repository revision, release page or tag, guiding the automation through its directory structure.

To provide your tasks and pipelines to this published version of this repository (the `p` branch), please follow this [workflow](./docs/workflow-provide-your-tekton-resources.md).

[artifactHub]: https://github.com/artifacthub.io
[artifactHubRedHatExperimentalPipelines]: https://artifacthub.io/packages/search?repo=redhat-tekton-experimental-pipelines
[artifactHubRedHatExperimentalTasks]: https://artifacthub.io/packages/search?repo=redhat-tekton-experimental-tasks
[artifactHubRedHatPipelines]: https://artifacthub.io/packages/search?repo=redhat-tekton-pipelines&page=1
[artifactHubRedHatTasks]: https://artifacthub.io/packages/search?repo=redhat-tekton-tasks&page=1
[openshiftPipelines]: https://docs.openshift.com/pipelines/1.12/about/about-pipelines.html
[openshiftPipelinesConfig]: https://docs.openshift.com/pipelines/1.12/create/remote-pipelines-tasks-resolvers.html#resolver-git-config-anon_remote-pipelines-tasks-resolvers
[pBranch]: https://github.com/openshift-pipelines/tektoncd-catalog/tree/p
[pipelineAsCode]: https://pipelinesascode.com/
[pipelineAsCodeRepository]: https://github.com/openshift-pipelines/pipelines-as-code/blob/main/docs/content/docs/guide/repositorycrd.md
[tektonGitResolver]: https://tekton.dev/docs/pipelines/git-resolver/
[tektonGitResolverPipeline]:https://tekton.dev/docs/pipelines/git-resolver/#pipeline-resolution
[tektonResolvers]: https://tekton.dev/docs/pipelines/resolution-getting-started/
[TEP0003]: https://github.com/tektoncd/community/blob/main/teps/0003-tekton-catalog-organization.md
