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
