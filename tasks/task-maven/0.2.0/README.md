# `maven` Tekton Task

The `maven` Task can be used to run a Maven goal on a simple or multi-module maven project.

## Workspaces

### `source`

The `source` workspace is required. It contains the source of the "maven" project to build. It should contain a `pom.xml`.

## `server_secret` (optional)

The `server_secret` is optional. It should contain two *files* : `username` and `password`. It is possible to bind a `ConfigMap` or a `Secret` to this workspace ; in that case, the `ConfigMap` or `Secret` should have a `username` and a `password` key.

## Parameters

| Parameter          | Type     | Default    | Description                                                                        |
|:-------------------|:---------|:-----------|:-----------------------------------------------------------------------------------|
| `GOALS`            | `string` | `package`  | The `maven` goal(s) to run                                                         |
| `MAVEN_MIRROR_URL` | `string` | "" (empty) | The maven repository mirror URL to use                                             |
| `SUBDIRECTORY`     | `string` | `.`        | The subdirectory of the `source` workspace on which we want to execute maven goals |

