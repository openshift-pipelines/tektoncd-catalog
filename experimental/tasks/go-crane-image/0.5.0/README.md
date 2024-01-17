# go-crane-image

Build an oci using go and crane.

- **Go** 1.20.
- **Crane** 0.17.x.
- The image(s) are based of Alpine.


## Workspaces

| Workspace      | Optional | Description                                            |
|:---------------|:--------:|:-------------------------------------------------------|
| `source`       | `false`  | The go source to build                                 |
| `dockerconfig` | `true`   | Includes a docker `config.json` or `.dockerconfigjson` |

## Params

| Param     | Type     | Default                                                      | Description                                                                                                                                    |
|:----------|:--------:|:-------------------------------------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------|
| `app`     | `string` | (required)                                                   | The name of the "application" to build. This will have an impact on the binary and possibly the image reference                                |
| `package` | `string` | `.`                                                          | The package to build. It needs to be a package `main` that compiles into a binary. The default value is `.`, usual value can be `./cmd/{name}` |
| `image`   | `object` | `{ base="", envs="", labels="", push="true", tag="latest" }` | The image specific options such as prefix, labels, env, …                                                                                      |
| `go`      | `object` | `{ GOOS="", CGO_ENABLED="0", GOARCH="", GOFLAGS="-v" }`      | Golang options, such as flags, version, …                                                                                                      |

## Results

| Result         | Description                     |
|:---------------|:--------------------------------|
| `IMAGE_DIGEST` | Digest of the image just built. |
| `IMAGE_URL`    | URL of the image just built.    |
