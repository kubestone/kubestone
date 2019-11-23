title: Kubestone - Fio: Flexible I/O performance benchmark

# Fio - Flexible I/O tester

!!! quote
    fio is a tool that will spawn a number of threads or processes doing a particular type of I/O action as specified by the user. The typical use of fio is to write a job file matching the I/O load one wants to simulate.

With the [fio](https://fio.readthedocs.io/en/latest/fio_doc.html) benchmark you can measure the I/O performance of the disks used in your Kubernetes cluster.



## Mode of operation

Kubestone generates a Kubernetes Job from each fio CR that will run a single pod with the defined fio job.

When `customJobFiles` are specified in the CR a ConfigMap will be created to hold the content of the job files. The entries in the ConfigMap named using the following pattern: `customJobN`, where N is the item in the customJobFiles list.

`Volume` defines the volume to use for benchmarking. 
`Volume.VolumeSource` provides way to mount already existing PVCs, HostPath, EmptyDir (and others) to the benchmark. 

When `Volume.PersistentVolumeClaimSpec` is defined (and `Volume.VolumeSource.PersistentVolumeClaim.ClaimName` set to 'GENERATED') a new PVC will be created for the benchmark. Note: The created volume is not freed up or removed after the benchmark run.


## Example configuration
You can find [configuration examples](https://github.com/xridge/kubestone/tree/master/config/samples/fio) in the GitHub repository.



## Sample benchmark
To run a sample benchmark with PVC mode, the following command can be used:
```bash
$ kustomize build github.com/xridge/kubestone/config/samples/fio/overlays/pvc | kubectl create --namespace kubestone -f -
```

Please refer to the [quickstart guide](../quickstart.md) for further details.




## Fio Configuration

The complete documentation of fio CR can be found in the [API Docs](../apidocs.md#perf.kubestone.xridge.io/v1alpha1.FioSpec).




## Docker Image

[Docker Image for Fio](https://hub.docker.com/r/xridge/fio) is provided via [xridge's fio repository](https://github.com/xridge/fio-docker).



## Legal

Fio is licensed as GPLv2.
