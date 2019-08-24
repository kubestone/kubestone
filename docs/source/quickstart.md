# Kubestone Quickstart Guide



Welcome to Kubestone, the benchmarking operator for Kubernetes.

This guide will walk you through on the one the installation process and will show you how to create a benchmark run on your cluster.



## Installation

**Requirements**

- [Kubernetes](https://kubernetes.io/) v1.14 (or newer)

- [Kustomize](https://kustomize.io/) v3.10

- Cluster admin privileges

  

Deploy Kubestone to `kubestone-system` namespace with the following command:

```bash
$ kustomize build github.com/xridge/kubestone/config/default | kubectl apply -f -
```



Once deployed, Kubestone will listen for Custom Resources created with the `kubestone.xridge.io` group.

If you wish to uninstall the operator please follow the [uninstall guide](uninstall.md).



## Benchmarking

Benchmarks can be executed via Kubestone by creating Custom Resources in your cluster. 

### Namespace

It is recommended to create a dedicated namespace for benchmarking.

```bash
$ kubectl create namespace kubestone
```

After the namespace is created you can use it to post a benchmark request to the cluster. 

The resulting  benchmark executions will reside in this namespace. 



### Custom Resource rendering

We will be using [kustomize](https://kustomize.io/) to render the Custom Resource from the [github repository](https://github.com/xridge/kubestone/tree/master/config/samples/fio/).

Kustomize takes a [base yaml](https://github.com/xridge/kubestone/blob/master/config/samples/fio/base/fio_cr.yaml) and patches with an [overlay file](https://github.com/xridge/kubestone/blob/master/config/samples/fio/overlays/pvc/patch.yaml) to render the final yaml file, which describes the benchmark. 

```bash
$ kustomize build github.com/xridge/kubestone/config/samples/fio/overlays/pvc
```

*Note: Kustomize has [extensive documentation](https://github.com/kubernetes-sigs/kustomize/tree/master/docs) about it's operation. It is advised to have a basic understanding on how Kustomize works as Kubestone uses it's features for both deployment and benchmark execution.*



The Custom Resource (rendered yaml) looks as follows:

```yaml
apiVersion: perf.kubestone.xridge.io/v1alpha1
kind: Fio
metadata:
  name: fio-sample
spec:
  cmdLineArgs: --name=randwrite --iodepth=1 --rw=randwrite --bs=4m --size=256M
  image:
    name: xridge/fio:3.13
  persistentVolumeClaim:
    accessModes:
    - ReadWriteOnce
    size: 5G
```



When we create this resource in Kubernetes, the operator interprets it and creates the associated benchmark. The fields of the Custom Resource controls how the benchmark will be executed:

- `metadata.name`: Identifies the Custom Resource. Later, this can be used to query or delete the benchmark in the cluster.
- `cmdLineArgs`: Arguments passed to the benchmark. This case we are providing the arguments to fio, a filesystem benchmark. It instructs the benchmark to execute a random write test with 4Mb of block size  with an overall transfer size of 256 MB.
- `image.name`: Describes the Docker Image of the benchmark. In case of [Fio](https://fio.readthedocs.io/) we are using [xridge's fio Docker Image](https://cloud.docker.com/u/xridge/repository/docker/xridge/fio), which is built from [this repository](https://github.com/xridge/fio-docker/).
- `persistentVolumeClaim`: Given that Fio is a disk benchmark we can set a PersistentVolumeClaim for the benchmark to be executed. The above setup instructs Kubernetes to take 5GB of space from the default StorageClass and use it for the benchmark.



### Running the benchmark

Now, as we understand the definition of the benchmark we can try to execute it.

*Note: Make sure you installed the operator and have it running, before executing this step.*

```bash
$ kustomize build github.com/xridge/kubestone/config/samples/fio/overlays/pvc | kubectl create --namespace kubestone -f -
```

Since we pipe the output of the `kustomize build` command into `kubectl create` it will create the object in our Kubernetes cluster. 



The resulting object can be queried using the object's type (`fio`)  and it's name (`fio-sample`):

```bash
$ kubectl describe --namespace kubestone fio fio-sample
Name:         fio-sample
Namespace:    kubestone
Labels:       <none>
Annotations:  <none>
API Version:  perf.kubestone.xridge.io/v1alpha1
Kind:         Fio
Metadata:
  Creation Timestamp:  2019-08-24T17:22:21Z
  Generation:          1
  Resource Version:    25337705
  Self Link:           /apis/perf.kubestone.xridge.io/v1alpha1/namespaces/kubestone/fios/fio-sample
  UID:                 bb18f3d3-c693-11e9-8071-4439c4920abc
Spec:
  Cmd Line Args:  --name=randwrite --iodepth=1 --rw=randwrite --bs=4m --size=256M
  Image:
    Name:         xridge/fio:3.13
    Pull Policy:  Always
  Persistent Volume Claim:
    Access Modes:
      ReadWriteOnce
    Size:  5Gi
Status:
  Completed:  true
  Running:    false
Events:
  Type    Reason           Age   From       Message
  ----    ------           ----  ----       -------
  Normal  CreateSucceeded  25s   kubestone  Created /api/v1/namespaces/kubestone/configmaps/fio-sample
  Normal  CreateSucceeded  25s   kubestone  Created /api/v1/namespaces/kubestone/persistentvolumeclaims/fio-sample
  Normal  CreateSucceeded  25s   kubestone  Created /apis/batch/v1/namespaces/kubestone/jobs/fio-sample
```



As the `Events` section shows, Kubestone has created a `ConfigMap`, a `PersistentVolumeClaim` and a` Job` for the provided Custom Resource. The `Status` field tells us that the benchmark has completed.



### Inspecting the benchmark

The created objects related to the benchmark can be listed using `kubectl` command:

```
$ kubectl get all --namespace kubestone
```

