# Fio - Flexible I/O tester

!!! quote
    fio is a tool that will spawn a number of threads or processes doing a particular type of I/O action as specified by the user.  The typical use of  fio  is  to  write  a  job  file matching the I/O load one wants to simulate. 

With the [fio](https://fio.readthedocs.io/en/latest/fio_doc.html) benchmark you can measure the I/O performance of the disks used in your Kubernetes cluster. 



## Mode of operation

Kubestone generates a Kubernetes Job from each fio CR that will run a single pod with the defined fio job.

When `customJobFiles` are specified in the CR a ConfigMap will be created to held the content of the job files. The entries in the ConfigMap named using the following pattern: `customJobN`, where N is the item in the customJobFiles list.

If `PersistentVolumeClaim` is provided in the CR the respective PVC will be used for benchmarking.

!!! warning
    If the PersistentVolumeClaim is not specified, Docker's layered fs performance will be measured

## Sample benchmark
To run a sample benchmark with PVC mode, the following command can be used:
```
$ kustomize build github.com/xridge/kubestone/config/samples/fio/overlays/pvc | kubectl create --namespace kubestone -f -
```

Please refer to the [quickstart guide](quickstart.md) for further details.




## Reference configuration
You can find [configuration examples](https://github.com/xridge/kubestone/tree/master/config/samples/fio) in the GitHub repository.



## Further documentation

The complete documentation of fio (and other benchmarks') CR can be found in the API Docs.



## Docker Image

[Docker Image for Fio](https://hub.docker.com/r/xridge/fio) is provided via [xridge's fio repository](https://hub.docker.com/r/xridge/fio).



## Legal

Fio is licensed as GPL v2. 