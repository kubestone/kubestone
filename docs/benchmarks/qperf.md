title: Kubestone - Qperf: RDMA and IP performance benchmark

# Qperf - Measure RDMA and IP performance

!!! quote
    Qperf measures bandwidth and latency between two nodes. It can work over
    TCP/IP as well as the RDMA transports. On one of the nodes, qperf is
    typically run with no arguments designating it the server node. One may
    then run qperf on a client node to obtain measurements such as bandwidth,
    latency and cpu utilization.

    In its most basic form, qperf is run on one node in server mode by invoking
    it with no arguments. On the other node, it is run with two arguments: the
    name of the server node followed by the name of the test.

With the [qperf](https://github.com/linux-rdma/qperf) benchmark, you can
measure the I/O performance of the network hardware and stack used in your
Kubernetes cluster.



## Mode of operation

As qperf requires a server and a client the controller creates the following objects during benchmark:
* Server Deployment
* Server Service
* Client Pod

At the first step, the Server Deployment and Service are created. Once both
becomes available, the Client Pod is created to execute the benchmark. Once the
benchmark is completed (regardless of it's success), the server deployment and
service is deleted from Kubernetes.

In order to avoid measuring loopback performance, it is advised that you set
the affinity and anti-affinity scheduling primitives for the benchmark. The
provided sample benchmark shows how to avoid executing the client and the
server on the same machine. For further documentation please refer to
Kubernetes' [respective documentation](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/).



## Example configuration

You can find [configuration example](https://github.com/xridge/kubestone/blob/master/config/samples/perf_v1alpha1_qperf.yaml) in the GitHub repository.



## Sample benchmark
```bash
$ kubectl create --namespace kubestone -f https://raw.githubusercontent.com/xridge/kubestone/master/config/samples/perf_v1alpha1_qperf.yaml
```


Please refer to the [quickstart guide](../quickstart.md) for details on generic principles and setup of Kubestone.




## Qperf Configuration

The complete documentation of qperf CR can be found in the [API Docs](../apidocs.md#perf.kubestone.xridge.io/v1alpha1.QperfSpec).



## Docker Image

[Docker Image for qperf](https://hub.docker.com/r/xridge/qperf) is provided via [xridge's qperf-docker repository](https://github.com/xridge/qperf-docker).



## Legal

Qperf is licensed as GPLv2.
