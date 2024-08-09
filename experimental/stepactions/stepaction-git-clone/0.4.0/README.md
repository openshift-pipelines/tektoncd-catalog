# `git-clone`

**Note: this StepAction is only compatible with Tekton Pipelines versions 0.54.0 and greater!**

**Note: this StepAction is not backwards compatible with the previous versions as it is now run as a non-root user!**

This `StepAction` has two required inputs:

1. The URL of a git repo to clone provided with the `URL` param.
2. A path called `OUTPUT_PATH` which should be provided by a `Workspace/Volume` in a Task.

The `git-clone` `StepAction` will clone a repo from the provided `URL` into the
`OUTPUT_PATH` directory. By default the repo will be cloned into the root of
your directory. You can clone into a subdirectory by setting this `StepAction`'s
`SUBDIRECTORY` param. If the directory where the repo will be cloned is
already populated then by default the contents will be deleted before the
clone takes place. This behaviour can be disabled by setting the
`DELETE_EXISTING` param to `"false"`.

**Note**: The `git-clone` StepAction is run as nonroot. The files cloned on to the `OUTPUT_PATH`
directory will end up owned by user 65532.

## Workspace-Paths

**Note**: `StepAction`s do not accept `Workspaces`. They only need the path to a directory 
which can be passed in via a `param`. The directory can be provided by the `Task` via a 
`Workspace`.

**Note**: This stepaction is executed in a Task as a non-root user with UID 65532 and GID 65532.
Generally, the default permissions for storage volumes are configured for the
root user. To make the volumes accessible by the non-root user, you will need
to either configure the permissions manually or set the `fsGroup` field under
`PodSecurityContext` in your TaskRun or PipelineRun.

An example PipelineRun will look like:
```yaml
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata:
  generateName: git-clone-
spec:
  pipelineRef:
    name: git-clone-pipeline
  podTemplate:
    securityContext:
      fsGroup: 65532
...
...
```

An example TaskRun will look like:
```yaml
apiVersion: tekton.dev/v1beta1
kind: TaskRun
metadata:
  name: taskrun
spec:
  taskRef:
    name: git-clone
  podTemplate:
    securityContext:
      fsGroup: 65532
...
...
```

* **output**: An example workspace for this StepAction to fetch the git repository in to.
  **Note**: The underlying StepAction only needs the path which can be provided via 
  the param as follows: 
  ```yaml
  params:
    - name: OUTPUT_PATH
      value: $(workspaces.output.path)
  ```
* **ssh-directory**: An optional workspace to provide SSH credentials. At
  minimum this should include a private key but can also include other common
  files from `.ssh` including `config` and `known_hosts`. It is **strongly**
  recommended that this workspace be bound to a Kubernetes `Secret`.
  **Note**: The underlying StepAction only needs the path which can be provided via 
  the param as follows: 
  ```yaml
  params:
    - name: SSH_DIRECTORY_PATH
      value: $(workspaces.ssh-directory.path)
  ```

* **ssl-ca-directory**: An optional workspace to provide custom CA certificates.
  Like the /etc/ssl/certs path this directory can have any pem or cert files,
  this uses libcurl ssl capath directive. See this SO answer here
  https://stackoverflow.com/a/9880236 on how it works.
  **Note**: The underlying StepAction only needs the path which can be provided via 
  the param as follows: 
  ```yaml
  params:
    - name: SSL_CA_DIRECTORY_PATH
      value: $(workspaces.ssa-ca-directory.path)
  ```

* **basic-auth**: An optional workspace containing `.gitconfig` and
  `.git-credentials` files. This allows username/password/access token to be
  provided for basic auth.

  It is **strongly** recommended that this workspace be bound to a Kubernetes
  `Secret`. For details on the correct format of the files in this Workspace
  see [Using basic-auth Credentials](#using-basic-auth-credentials) below.

  **Note**: Settings provided as part of a `.gitconfig` file can affect the
  execution of `git` in ways that conflict with the parameters of this StepAction.
  For example, specifying proxy settings in `.gitconfig` could conflict with
  the `HTTP_PROXY` and `HTTPS_PROXY` parameters this StepAction provides. Nothing
  prevents you setting these parameters but it is not advised.
  
  **Note**: The underlying StepAction only needs the path which can be provided via 
  the param as follows:
  ```yaml
  params:
    - name: BASIC_AUTH_PATH
      value: $(workspaces.basic-auth.path)
  ```

## Parameters

* **OUTPUT_PATH**: The git repo will be cloned onto this path
* **SSH_DIRECTORY_PATH**: A .ssh directory with private key, known_hosts, config, etc.
* **BASIC_AUTH_PATH**:  A directory path containing a .gitconfig and .git-credentials file.
* **SSL_CA_DIRECTORY_PATH**: A directory containing CA certificates, this will be used by Git to verify the peer with when fetching or pushing over HTTPS.
* **URL**: Repository URL to clone from. (_required_)
* **REVISION**: Revision to checkout. (branch, tag, sha, ref, etc...) (_default_: "")
* **REFSPEC**: Refspec to fetch before checking out revision. (_default_:"")
* **SUBMODULES**: Initialize and fetch git submodules. (_default_: true)
* **DEPTH**: Perform a shallow clone, fetching only the most recent N commits. (_default_: 1)
* **SSL_VERIFY**: Set the `http.sslVerify` global git config. Setting this to `false` is not advised unless you are sure that you trust your git remote. (_default_: true)
* **CRT_FILENAME**: If `sslVerify` is **true** and `ssl-ca-directory` workspace is given then set `crtFileName` if mounted file name is different than `ca-bundle.crt`. (_default_: "ca-bundle.crt")
* **SUBDIRECTORY**: Subdirectory inside the `output` workspace to clone the repo into. (_default:_ "")
* **DELETE_EXISTING**: Clean out the contents of the destination directory if it already exists before cloning. (_default_: true)
* **HTTP_PROXY**: HTTP proxy server for non-SSL requests. (_default_: "")
* **HTTPS_PROXY**: HTTPS proxy server for SSL requests. (_default_: "")
* **NO_PROXY**: Opt out of proxying HTTP/HTTPS requests. (_default_: "")
* **VERBOSE**: Log the commands that are executed during `git-clone`'s operation. (_default_: true)
* **SPARSE_CHECKOUT_DIRECTORIES**: Which directories to match or exclude when performing a sparse checkout (_default_: "")
* **USER_HOME**: The user's home directory. (_default_: "/tekton/home")

## Results

* **COMMIT**: The precise commit SHA that was fetched by this StepAction
* **URL**: The precise URL that was fetched by this StepAction
* **COMMITTER_DATE**: The epoch timestamp of the commit that was fetched by this StepAction

## Platforms

The StepAction can be run on `linux/amd64`, `linux/s390x`, `linux/arm64`, and `linux/ppc64le` platforms.

## Usage

If the `REVISION` is not provided in the param of the taskrun
then it will auto-detect the branch as specified by the `default`
in the respective git repository.

The following pipelines demonstrate usage of the git-clone StepAction:

## Cloning Private Repositories

This StepAction supports fetching private repositories. There are three ways to
authenticate:

1. The simplest approach is to bind an `ssh-directory` workspace to this
StepAction. The workspace should contain private keys (e.g. `id_rsa`), `config`
and `known_hosts` files - anything you need to interact with your git remote
via SSH. It's **strongly** recommended that you use Kubernetes `Secrets` to
hold your credentials and bind to this workspace.

    In a TaskRun that would look something like this:

    ```yaml
    kind: TaskRun
    spec:
      workspaces:
      - name: ssh-directory
        secret:
          secretName: my-ssh-credentials
    ```

    And in a Pipeline and PipelineRun it would look like this:

    ```yaml
    kind: Pipeline
    spec:
      workspaces:
      - name: ssh-creds
      # ...
      tasks:
      - name: fetch-source
        taskSpec:
          workspaces:
            - name: ssh-creds
          steps:
          - ref:
              name: git-clone
            params:
              - name: SSH_DIRECTORY_PATH
                value: $(workspaces.ssh-creds.path)
        workspaces:
          -name: ssh-creds
      # ...
    ---
    kind: PipelineRun
    spec:
      workspaces:
      - name: ssh-creds
        secret:
          secretName: my-ssh-credentials
      # ...
    ```

    The `Secret` would appear the same in both cases - structured like a `.ssh`
    directory:

    ```yaml
    kind: Secret
    apiVersion: v1
    metadata:
      name: my-ssh-credentials
    data:
      id_rsa: # ... base64-encoded private key ...
      known_hosts: # ... base64-encoded known_hosts file ...
      config: # ... base64-encoded ssh config file ...
    ```

    Including `known_hosts` is optional but strongly recommended. Without it
    the `git-clone` StepAction will blindly accept the remote server's identity.

2. Use Tekton Pipelines' built-in credentials support as [documented in
Pipelines' auth.md](https://github.com/tektoncd/pipeline/blob/master/docs/auth.md).

3. Another approach is to bind an `ssl-ca-directory` workspace to this
StepAction. The workspace should contain crt keys (e.g. `ca-bundle.crt`)files - anything you need to interact with your git remote
via custom CA . It's **strongly** recommended that you use Kubernetes `Secrets` to
hold your credentials and bind to this workspace.

    In a TaskRun that would look something like this:

    ```yaml
    kind: TaskRun
    spec:
      workspaces:
      - name: ssl-ca-directory
        secret:
          secretName: my-ssl-credentials
    ```

    And in a Pipeline and PipelineRun it would look like this:

    ```yaml
    kind: Pipeline
    spec:
      workspaces:
      - name: ssl-creds
      # ...
      tasks:
      - name: fetch-source
        workspaces:
        - name: ssl-creds
        taskSpec:
          workspaces:
            - name: ssl-creds
          steps:
          - ref:
              name: git-clone
            params:
              - name: SSL_CA_DIRECTORY_PATH
                value: $(workspaces.ssl-creds.path)
      # ...
    ---
    kind: PipelineRun
    spec:
      workspaces:
      - name: ssl-creds
        secret:
          secretName: my-ssl-credentials
      # ...
    ```

    The `Secret` would appear like below:

    ```yaml
    kind: Secret
    apiVersion: v1
    metadata:
      name: my-ssl-credentials
    data:
      ca-bundle.crt: # ... base64-encoded crt ...  # If key/filename is other than ca-bundle.crt then set crtFileName param as explained under Parameters section
    ```

## Using basic-auth Credentials

**Note**: It is strongly advised that you use `ssh` credentials when the option
is available to you before using basic auth. You can generate a short
lived token from WebVCS platforms (Github, Gitlab, Bitbucket etc..) to be used
as a password and generally be able to use `git` as the username.
On bitbucket server the token may have a / into it so you would need
to urlquote them before in the `Secret`, see this stackoverflow answer :

https://stackoverflow.com/a/24719496 

To support basic-auth this StepAction exposes an optional `basic-auth` Workspace.
The bound Workspace must contain a `.gitconfig` and `.git-credentials` file.
Any other files on this Workspace are ignored. A typical `Secret` containing
these credentials looks as follows:

```yaml
kind: Secret
apiVersion: v1
metadata:
  name: my-basic-auth-secret
type: Opaque
stringData:
  .gitconfig: |
    [credential "https://<hostname>"]
      helper = store
  .git-credentials: |
    https://<user>:<pass>@<hostname>
```
