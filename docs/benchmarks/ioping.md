title: Kubestone - ioping: Storage Latency benchmark

# ioping - Storage Latency benchmark

!!! quote
    A tool to monitor I/O latency in real time. It shows disk latency in the same way as ping shows network latency.

With [ioping](https://github.com/koct9i/ioping) benchmark you can measure the latency of the storage I/O subsystem in your Kubernetes cluster.



## Mode of operation

Kubestone generates a Kubernetes Job from each ioping CR.

`Volume` defines the volume to use for benchmarking. 
`Volume.VolumeSource` provides way to mount already existing PVCs, HostPath, EmptyDir (and others) to the benchmark. 

When `Volume.PersistentVolumeClaimSpec` is defined (and `Volume.VolumeSource.PersistentVolumeClaim.ClaimName` set to 'GENERATED') a new PVC will be created for the benchmark. Note: The created volume is not freed up or removed after the benchmark run.



## Example configuration
You can find [configuration example](https://github.com/xridge/kubestone/blob/master/config/samples/perf_v1alpha1_ioping.yaml) in the GitHub repository.


## Sample benchmark
To run a sample benchmark with EmptyDir, the following command can be used:
```bash
$ kubectl create --namespace kubestone -f https://raw.githubusercontent.com/xridge/kubestone/master/config/samples/perf_v1alpha1_ioping.yaml
```

Please refer to the [quickstart guide](../quickstart.md) for further details.




## ioping configuration

The complete documentation of ioping CR can be found in the [API Docs](../apidocs.md#perf.kubestone.xridge.io/v1alpha1.IopingSpec).




## Docker Image

[Docker Image for ioping](https://hub.docker.com/r/xridge/ioping) is provided via [xridge's ioping repository](https://github.com/xridge/ioping-docker).



## Legal

ioping is licensed as GPLv3.
