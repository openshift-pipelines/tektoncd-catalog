# Design

From repository `task-git` to `tektoncd-catalog` `p` branch.

- In `task-git` repository
	- I have a "configuration file" (`catalog.yaml` ?) that lists Tasks and Pipelines from the repository that I want to release
		- This file can be generated, updated, … 
	- When I do the release, I want to issue one command (`catalog-cd release`)
		- it mutates the resources to add the version annotation
		- it generates the final `catalog.yaml` with hash, digest, signature, …
		- it packages the tasks and pipelines in a `tekton-resources.tar.gz` tarball (with READMEs for documentation)
		- it (optionally) create, push the tag, create a GitHub release and attach content to it
- In `tektoncd-catalog` repository
	- `task-git` is configured in the `externals.yaml` configuration file
	- A schedule action (each hours ?) does the following, for each entry in `externals.yaml`
		- List releases and filter those that have a `catalog.yaml`
		- Fetch the `catalog.yaml` and the `tekton-resources.tar.gz`
		- Extract the tarball content and merge it with the current catalog available in the `p` branch
		- Creates a pull-request to update it
	- The pull request checks includes
		- Lint the resources
		
What is describe above is *required* for the internal launch.

What is missing from here:
- Attestation, SBOM, signature, …
- How to validate the task is well tested (so that Red Hat can support it)

---

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
