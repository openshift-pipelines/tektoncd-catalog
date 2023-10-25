Red Hat Tekton Ecosystem Team
-----------------------------

The Tekton Ecosystem team aims to provide well written Tekton resources like Tasks, Pipeline, Triggers, and any other components used in combination with [`tektoncd/pipeline`][tektonPipeline].

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

[tektonPipeline]: https://github.com/tektoncd/pipeline
