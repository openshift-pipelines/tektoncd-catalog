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
