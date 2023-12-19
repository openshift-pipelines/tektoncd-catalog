<p align="center">
	<a alt="Test Workflow" href="https://github.com/openshift-pipelines/task-git/actions/workflows/test.yaml">
		<img alt="GitHub Test Workflow Status" src="https://img.shields.io/github/actions/workflow/status/openshift-pipelines/task-git/test.yaml?label=test">
	</a>
	<a alt="Latest Release" href="https://github.com/openshift-pipelines/task-git/releases/latest">
		<img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/openshift-pipelines/task-git">
	</a>
</p>

`git` Tekton Task
-----------------

The `git` Task will clone and prepare a Git repository on the `output` Workspace. By default the Git repository (`URL`) will be cloned on the root of the Workspace, but you can alter the clone location using the `SUBDIRECTORY` parameter, and [other settings](#parameters).

A quick usage example is:

```yaml
---
apiVersion: tekton.dev/v1beta1
kind: TaskRun
metadata: {}
spec:
  taskRef:
    name: git
  params:
    - name: URL
      value: https://github.com/openshift-pipelines/task-git.git
  workspaces:
    - name: output
      emptyDir: {}
```

Please consider the [Workspaces](#workspaces), [Parameters](#parameters) and [Results](#results) described below.

# Workspaces

A single Workspace is required for this Task, while the optional Workspaces will allow advanced Git configuration and authentication.

## `output`

The `output` is a required Workspace, represents the primary location where the Git repository data will be stored, and later on, this Workspace will be monted in subsequent Pipeline's Task.

Knowing the Workspace data will be employed on other Tasks, the recommended approach is using a [persistent volume][tektonPVC], for instance a [`PersistentVolumeClaim` (PVC)][k8sPVC].

## Authentication Workspaces

The recommended approach to authentication is using the [default mechanisms supported by Tekton Pipeline][tektonAuthentication], please consider it as your first option.

More advanced use-cases may require different methods of interacting with private repositories, the following Workspaces are meant to support advanced Git configuration and authentication.

### `basic-auth` (HTTP/SSH)

The `basic-auth` is a optional Workspace to provide Git credentials and configuration.

The following Workspace files (items) are shared with Git before cloning the repository, the Task copies the files to the Git user home directory, configured by the parameter `USER_HOME`.

| Workspace File     | Required | Description                            |
| :----------------- | :------: | :------------------------------------- |
| `.git-credentials` |  `true`  | [Git credentials file][gitCredentials] |
| `.gitconfig`       |  `true`  | [Git configuration file][gitConfig]    |

Typically, this type of data is stored as a Kubernetes Secret. For example:

```bash
kubectl create secret generic basic-auth-ex \
	--from-file=".git-credentials=${HOME}/.git-credentials" \
	--from-file=".gitconfig=${HOME}/.gitconfig"
```

Then, you can [reference the Secret][tektonWorkspaceSecret] as the `basic-auth` Workspace.

### `ssh-directory` (SSH)

The `ssh-directory` is a optional Workspace, meant to store the files commonly found on a [`~/.ssh` directory][dotSSHDirectory], when informed, the whole directory will be copied into the Git's home (configured by the parameter `USER_HOME`).

During the `prepare` step you can see the details about what's being copied, please consider the output log snippet below. For more verbose logging set the peramater `VERBOSE` to `true`.

```
---> Phase: Copying '.ssh' from ssh-directory workspace ('/workspaces/ssh-directory')...
'/workspaces/ssh-directory' -> '/home/git/.ssh'
'/workspaces/ssh-directory/config' -> '/home/git/.ssh/config'
mode of '/home/git/.ssh' changed from 0755 (rwxr-xr-x) to 0700 (rwx------)
mode of '/home/git/.ssh/config' changed from 0644 (rw-r--r--) to 0400 (r--------)

```

It's recommended storing this type of data as a Kubernetes Secret, like the following example:

```bash
kubectl create secret generic ssh-directory-ex \
	--from-file="config=${HOME}/.ssh/config" \
	--from-file="authorized_keys=${HOME}/.ssh/authorized_keys"
```

Then, you can [reference the Secret][tektonWorkspaceSecret] as the `ssh-directory` Workspace.

### `ssl-ca-directory` (mTLS)

The `ssl-ca-directory` is a optional Workspace to store a additional [Certificate Authority (CA)][tlsCA] bundles, commonly `.pem` or `.crt` files. The exact bundle file name is defined by the parameter `CRT_FILENAME`.

Before running the Git clone command, the [`GIT_SSL_CAINFO` environment variable][gitSSLCAInfo] is exported with the full path to the `CRT_FILENAME` in the `ssl-ca-directory` Workspace.

You can observe the setting taking place on the beggining of the `git-clone` step:

```
phase 'Exporting custom CA certificate "GIT_SSL_CAINFO=/workspaces/ssl-ca-directory/ca-bundle.crt"'
```

This is a sensitive information and therefore it's recommended to store as a Kubernetes Secret, please consider the previous examples to create Secrets with the `--from-file` option.

Finally, you can [reference the Secret][tektonWorkspaceSecret] as the `ssl-ca-directory` Workspace.

# Parameters

The following parameters are supported by this Task.

| Parameter | Type  | Default | Description |
| :-------- | :---: | :------ | :---------- |
| `URL` | `string` | (required) | Git repository URL|
| `REVISION` | `string` | `main` | Revision to checkout, an branch, tag, sha, ref, etc... |
| `SUBMODULES` | `string` | `true` | Initialize and fetch Git submodules |
| `DEPTH` | `string` | `1` | [Number of commits to fetch][gitCloneDepath], a "shallow clone" is a single commit |
| `SSL_VERIFY` | `string` | `true` | Sets the global [`http.sslVerify`][gitHTTPSSLVerify] value, `false` is not advised unless you trust the remote repository |
| `CRT_FILENAME` | `string` | `ca-bundle.crt` | Certificate Authority (CA) bundle filename on the `ssl-ca-directory` Workspace. |
| `SUBDIRECTORY` | `string` | "" (empty) |  Relative path to the `output` Workspace where the repository will be cloned |
| `SPARSE_CHECKOUT_DIRECTORIES` | `string` | "" (empty) | List of directory patterns split by comma to perform ["sparse checkout"][gitSparseCheckout] |
| `DELETE_EXISTING` | `string` | `true` | Clean out the contents of the `output` Workspace before cloning the repository, if data exists. |
| `HTTP_PROXY` | `string` | "" (empty) | HTTP proxy server (non-TLS requests) |
| `HTTPS_PROXY` | `string` | "" (empty) | HTTPS proxy server (TLS requests) |
| `NO_PROXY` | `string` | "" (empty) | Opt out of proxying HTTP/HTTPS requests |
| `VERBOSE` | `string` | `false` | Log the commands executed |
| `USER_HOME` | `string` | `/home/git` | Absolute path to the Git user home directory |

# Results

The following results are produced by this Task.

| Name | Description |
| :--- | :---------- |
| `COMMIT` | The precise commit SHA digest cloned |
| `URL` | The precise repository URL |
| `COMMITTER_DATE` | The epoch timestamp of the commit cloned |

[dotSSHDirectory]: https://man.openbsd.org/sshd#FILES
[gitCloneDepath]: https://git-scm.com/docs/git-clone#Documentation/git-clone.txt---depthltdepthgt
[gitConfig]: https://git-scm.com/docs/git-config#FILES
[gitCredentials]: https://git-scm.com/docs/git-credential-store#Documentation/git-credential-store.txt-git-credentials
[gitHTTPSSLVerify]: https://git-scm.com/docs/git-config#Documentation/git-config.txt-httpsslVerify
[gitSparseCheckout]: https://git-scm.com/docs/git-sparse-checkout#_description
[gitSSLCAInfo]: https://git-scm.com/docs/git-config#Documentation/git-config.txt-httpsslCAInfo
[k8sPVC]: https://kubernetes.io/docs/concepts/storage/persistent-volumes/
[tektonAuthentication]: https://tekton.dev/docs/pipelines/auth/
[tektonPVC]: https://tekton.dev/docs/pipelines/workspaces/#using-persistentvolumeclaims-as-volumesource
[tektonWorkspaceSecret]: https://tekton.dev/docs/pipelines/workspaces/#secret
[tlsCA]: https://en.wikipedia.org/wiki/Certificate_authority
