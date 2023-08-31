Tekton Resources Linter
-----------------------

# Summary

This document describes the best practices enforced by the Tekton resources [linter][catalogLinter] on this catalog repository.

The objective is to guide Tekton resources creators/maintainers to adopt the same style and conventions throughtout the catalog, and provide a improved experience consuming the resources.

# Best Practices

The following describes the best practices for individual parts of Tekton Tasks and Pipelines.

## Red Hat Certified Images

In order to onboard on this catalog the Tekton resources container images, for Tasks and Pipelines, must be based on a certified Red Hat container image, please search for images on the [official catalog](https://catalog.redhat.com/software/containers/search).

## Tekton Resources Naming

Please consider the name convention for Tekton resources on the catalog, each applicable resource attribute is expected to contain the `.description` information.

Additionally, please consider other resources in the catalog in order to adopt common names for inputs. For instance, `buildah` and `s2i` tasks rely on the same base Workspace, Params and Results naming in order to improve re-usability; by knowing the expected terms used throughout the catalog resources results on a smoother overall experience.

### Workspaces

Workspaces must be lower case using dash (`-`) to split words, i.e.:

```yml
- name: build
- name: source
- name: test-stub
```

### Params

Params must be uppercase using underscore  (`_`) to split words, i.e:

```yml
- name: IMAGE
- name: PROCESS_TYPE
- name: USER_ID
```

### Results

Results follows the same convention than [Params](#params), therefore names must be uppercase using underscore (`_`) to split words, i.e:

```yml
- name: IMAGE_SHA256
- name: REGISTRY_URL
- name: REVISION
```

[catalogLinter]: https://github.com/openshift-pipelines/tektoncd-catalog/blob/main/internal/linter/linter.go
[redHatImagesCatalog]: https://catalog.redhat.com/software/containers/search
