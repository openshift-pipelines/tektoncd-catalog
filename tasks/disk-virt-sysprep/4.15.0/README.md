# Disk Virt Sysprep Task

This task uses [virt-sysprep](https://libguestfs.org/virt-sysprep.1.html) to run a sysprep script on a target pvc.

### Parameters

- **pvc**: PersistentVolumeClaim to run the the virt-sysprep script in. PVC should be in the same namespace as taskrun/pipelinerun.
- **sysprepCommands**: virt-sysprep commands in `--commands-from-file` format.
- **verbose**: Enable verbose mode and tracing of libguestfs API calls.
- **additionalOptions**: Additional options to pass to virt-sysprep.

### Usage

Task run using resolver:
```
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
    generateName: disk-virt-sysprep-taskrun-
spec:
    params:
    -   name: pvc
        value: example-pvc
    -   name: sysprepCommands
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
            value: disk-virt-sysprep
        -   name: version
            value: 4.15.0
        resolver: hub
```

#### Common Errors

- The input PVC disk should not be accessed by a running VM or other tools like virt-sysprep task concurrently.
The task will fail with a generic `...guestfs_launch failed...` message.
A verbose parameter can be set to true for more information.

### OS support

- Linux: full; all the sysprep commands work
- Windows: partial; only some sysprep commands work
