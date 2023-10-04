# tektoncd-catalog

Catalog of Tekton resources (Tasks, Pipelines, …) by Red Hat. 

This repository contains a catalog of `Task` resources (and someday
`Pipeline`s and other resources), which are designed to be reusable in many
pipelines, authored and supported by Red Hat.

Each `Task` is provided in a separate directory along with a README.md and a
Kubernetes manifest, so you can choose which `Task`s to install on your
cluster. A directory can hold one task and multiple versions.

The layout of this repository follows the of [TEP-003: Tekton catalog
organization](https://github.com/tektoncd/community/blob/main/teps/0003-tekton-catalog-organization.md).

# What is Tekton ecosystem team about

The Tekton Ecosystem team aims to provide well written Tekton payloads (Tasks, Pipeline,
Triggers, and any other element that can be used with `tektoncd/pipeline`).

+ Well written (easy to use, customizable), supported, Tasks and Pipelines.
  - Define a set of "support" levels (incubation, …)
+ Same for any other type we think we should provide to our customers
+ Allow Red Hat teams to own and maintain their own set of Tekton resources (for their project)
+ Publish, Document and "publicize" those resources
+ Provide easy way to compose and maintain tasks in Pipeline, in cluster, …
  - Provide tooling for this
  - Possibly provide `CustomTask` or `Resolvers` for it
  - Build on top of existing tools (`renovate`, …)

+ Long term
  - Certifications (for partners)
  - Red Hat catalog presense (catalog.redhat.com)
  - Possibly driving API "feature"/changes
    Because we will write a lot of task, use them, … we should be able to find gap or
    enhancements in the API, and propose them as TEPs (with data).

# Usage Examples


This section explains how to use the tasks supported in this repository with the help of various tools like [Tekton Resolvers](https://tekton.dev/docs/pipelines/resolution-getting-started/) as well as [Pipelines as Code](https://pipelinesascode.com/). 

## Using Tekton Resolvers

Make sure kubectl is installed, if not install it using this [link](https://kubernetes.io/docs/tasks/tools/).

To use our tasks, you can create a Pipeline Resource as follows:

```yaml
---
apiVersion: tekton.dev/v1beta1
kind: Pipeline
metadata:
  labels:
    name: example-pipeline
  name: example-pipeline
spec:
  params:
    # Customize the Params as needed by your chosen tasks
    - name: APP_NAME
      type: string
      default: example
    - name: IMAGE_PREFIX
      type: string
      default: "test"

  workspaces:
    - name: source

  tasks:
    # Add or Reference other tasks as needed
    - name: example-task
      taskRef:
        resolver: git
        params:
          - name: url
            value: https://github.com/openshift-pipelines/tektoncd-catalog.git
          - name: revision
            value: p
          - name: pathInRepo
            value: experimental/tasks/go-crane-image/v0.1.0/go-crane-image.yaml
      workspaces:
        - name: source
          workspace: source
      params:
        - name: app
          value: $(params.APP_NAME)
        - name: image
          value:
            prefix: $(params.IMAGE_PREFIX)
```

For this example we have used a PersistentVolumeClaim as follows:

```yaml
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  labels:
    name: test
  name: test
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 250Mi
```

Then use the following commands to apply & run the above Pipeline

- Create PVC resource: ```kubectl apply -f pvc.yaml```
- Create Pipeline: ```kubectl apply -f pipeline.yaml```
- PipelineRun: ```tkn pipeline start example-pipeline --workspace="name=source,claimName=test,subPath=source" --showlog```

To learn more about resolver, use this [link](https://tekton.dev/docs/pipelines/resolution-getting-started/). 

## Using Pipelines as Code

Make sure all the prerequisites are present for using pac (check that [here](https://pipelinesascode.com/docs/install/getting-started/)

To create a template PipelineRun resource, you can use the command: tkn pac generate
After this it'll create a PipelineRun template, which can then be customized similar to the previous example to incorporate our tasks

Some annotations to look for are: 
```
metadata:
  annotations:
    pipelinesascode.tekton.dev/on-event: "[push]"
    pipelinesascode.tekton.dev/on-target-branch: "[main]"
    pipelinesascode.tekton.dev/max-keep-runs: "5"
```

## Cloning the repo or Directly (Not Recommended)

You can also consume our resources by cloning this repository and manually creating the resources in the cluster as well

Helpful commands:
- `kubectl apply -f https://github.com/openshift-pipelines/tektoncd-catalog/blob/p/experimental/tasks/name-of-task/version/file.yaml`, replace "name-of-task/version/file" according to your required task
OR
- `git clone https://github.com/openshift-pipelines/tektoncd-catalog.git`
- `kubectl apply path-of-task.yaml` (replace path-of-task with relevant task's path)

After adding the Tasks to the cluster, you can use them as needed for other resources 
