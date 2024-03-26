# Disk Virt Customize Task

This task uses [virt-customize](https://libguestfs.org/virt-customize.1.html) to run a customize script on a target pvc.

### Parameters

- **pvc**: PersistentVolumeClaim to run the the virt-customize script in. PVC should be in the same namespace as taskrun/pipelinerun.
- **customizeCommands**: virt-customize commands in `--commands-from-file` format.
- **verbose**: Enable verbose mode and tracing of libguestfs API calls.
- **additionalOptions**: Additional options to pass to virt-customize.

### Usage

Task run using resolver:
```
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
    generateName: disk-virt-customize-taskrun-
spec:
    params:
    -   name: pvc
        value: example-pvc
    -   name: customizeCommands
        value: install make,ansible
    taskRef:
        params:
        -   name: catalog
            value: kubevirt-tekton-tasks
        -   name: type
            value: artifact
        -   name: kind
            value: task
        -   name: name
            value: disk-virt-customize
        -   name: version
            value: 4.15.0
        resolver: hub
```

#### Common Errors

- The input PVC disk should not be accessed by a running VM or other tools like virt-customize concurrently.
The task will fail with a generic `...guestfs_launch failed...` message.
Verbose parameter can be set to true for more information.

### OS support

- Linux: full; all the customize commands work
- Windows: partial; only some customize commands work
