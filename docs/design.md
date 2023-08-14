# Design

## Workflow

- `p` is where the indexed catalog is (where users can pull it from)
- The content of the `p` branch is generated from a set of configuration in the `main` branch
    - and the `catalog-cd` tool
    - it reads `externals.yaml` file (according to a schema) to know which repository to pull
    - for each repository, we go through each releases
        - for each releases, we fetch the `contract.yaml` to get the task and pipeline to pull
        - Based of the contract, we fetch the tekton resources (task, pipeline)
        - *we ignore failures at the moment*
    - We create the folder hierarchy, and submit a PR…
        - … with `lgtm`, `ok-to-test` and `approve` labels to be automatically merged
    - The PRs will run all the tests suites before it can be merged

## Tooling

The main tool used for this is `catalog-cd`. This is for generating the catalog as well as running lints and tests.

We are using GitHub Actions, in a daily fashion to keep generating it.

## Future

- Add support for specifying how to test the resource (command, using `catalog-cd`, …) in contracts
- Enhance `catalog-cd generate`
    - Generate documentation
- Instead of creating one big pull-request, we should create one per repository pulled so that we don’t block upgrading some tekton resources because others are not working (lint, tests, …)
- Have notifications (when failures, …)
