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


This section explains how to use the tasks supported in this repository with the help of various tools like Tekton Resolvers as well as Pipelines as Code. 

## Using Tekton Resolvers

Make sure kubectl is installed, if not install it using this [link](https://kubernetes.io/docs/tasks/tools/).

After that create a YAML file as follows:

```yaml
---
apiVersion: tekton.dev/v1beta1
kind: TaskRun
metadata:
  name: simple-taskrun-example
spec:
  workspaces:
    - name: source
      persistentVolumeClaim:
        claimName: your-claim-name
      subPath: source
  taskRef:
    resolver: git
    params:
      - name: url
        value: https://github.com/openshift-pipelines/tektoncd-catalog.git
      - name: revision
        value: p
      - name: pathInRepo
        value: experimental/tasks/go-crane-image/v0.1.0/go-crane-image.yaml
  params:
    - name: app
      value: example-task
    - name: image
      value:
        prefix: "add-custom-prefix"
```

Filename used in example is taskrun.yaml

Note that for this example we have used a PersistentVolumeClaim as follows:

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

Filename used in example is pvc.yaml

Then use the following commands to apply & run the above TaskRun

- Create PVC resource: kubectl apply -f pvc.yaml
- Create TaskRun: kubectl apply -f taskrun.yaml

Similarly you use Resolvers to create Pipelines & PipelineRuns as well. 

To learn more about resolver, use this [link](https://tekton.dev/docs/pipelines/resolution-getting-started/). 

## Using Pipelines as Code

WIP
