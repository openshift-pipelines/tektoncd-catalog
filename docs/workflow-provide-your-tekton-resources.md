# Provide task or pipeline to the catalog

This document is the "guide" to provide Tekton `Task` or `Pipeline` to the catalog, in an *almost* automated way.

As described in the [`design.md`](./design.md), we pull release content from each repository and populate it in the `p` branch, if it follows the [contract](./catalog.md). This means that in order to get your `Task` or `Pipeline` added to this repository you need to do the following:

- You need to have a repository with the `Task` and/or `Pipeline` you maintain. [`task-git`](https://github.com/openshift-pipelines/task-git) can serve as an example.
- You need to add your repository to the [`externals.yaml`](../externals.yaml) file (or the [experimental one](../experimental/externals.yaml) if it's a set of *more experimental resources*).
- You need to create tag and release for resources (`Task` and/or `Pipeline`) to be pull.
	- The release(s) need, *today*, two components:
		- a [`catalog.yaml`](./catalog.md) file
		- a `resources.tar.gz` file that contains your resources
	- In order to do the release, you can use our *experimental* [`catalog-cd`](../cmd/catalog-cd) tool, and its `release` sub-command to help, see [below](#release-with-catalogcd).

*Note: as of today, we only support GitHub releases, but the goal is to support more and more provides as we go*.

*Note: in the near future, we will also provide GitHub Action and Tekton Task and Pipeline to help you do the release.*

## Releasing with `catalog-cd`

You can use [`catalog-cd`](../cmd/catalog-cd) to prepare the release of your Tekton resources.

```shell
$ go run github.com/openshift-pipelines/tektoncd-catalog/cmd/catalog-cd release --output {path-for-the-release} --version="0.1.0" {paths-to-tasks-and-or-pipelines}
# e.g. catalog-cd release --output /tmp/release --version=0.1.0 task/go-crane-image pipeline/my-go-pipeline something/else
```

This command will generate the following:
- a [`catalog.yaml`](./catalog.md) file
- a `resources.tar.gz` file that contains your resources (organized by types)

Once you have those files, you can create and push a git tag, create a release and attach those files to it.

```shell
$ git tag 0.1.0
$ git push 0.1.0
$ gh release create 0.1.0 --generate-notes
$ gh release upload 0.1.0 release/catalog.yaml
$ gh release upload 0.1.0 release/resources.tar.gz
```

And that's all folks. If your repository is listed in the [`externals.yaml`](../externals.yaml) file, the new version will be picked up in the next few hours.
