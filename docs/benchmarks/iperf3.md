title: Kubestone - Iperf3: Network bandwidth performance benchmark

# Iperf3 - Network bandwidth benchmark

!!! quote
    iPerf3 is a tool for active measurements of the maximum achievable bandwidth on IP networks. It supports tuning of various parameters related to timing, buffers and protocols (TCP, UDP, SCTP with IPv4 and IPv6). 

With the [iperf3](https://iperf.fr/) benchmark, you can measure the I/O performance of the network hardware and stack used in your Kubernetes cluster. 



## Mode of operation

As iperf3 requires a server and a client the controller creates the following objects during benchmark:

- Server Deployment

- Server Service

- Client Pod

  

At the first step, the Server Deployment and Service are created. Once both becomes available, the Client Pod is created to execute the benchmark. Once the benchmark is completed (regardless of it's success), the server deployment and service is deleted from Kubernetes.

In order to avoid measuring loopback performance, it is advised that you set the affinity and anti-affinity scheduling primitives for the benchmark. The provided sample benchmark shows how to avoid executing the client and the server on the same machine. For further documentation please refer to Kubernetes' [respective documentation](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/).



## Example configuration

You can find [configuration example](https://github.com/xridge/kubestone/blob/master/config/samples/perf_v1alpha1_iperf3.yaml) in the GitHub repository.



## Sample benchmark
```bash
$ kubectl create --namespace kubestone -f https://raw.githubusercontent.com/xridge/kubestone/master/config/samples/perf_v1alpha1_iperf3.yaml
```


Please refer to the [quickstart guide](../quickstart.md) for details on generic principles and setup of Kubestone.




## IPerf3 Configuration

The complete documentation of iperf3 CR can be found in the [API Docs](../apidocs.md#perf.kubestone.xridge.io/v1alpha1.Iperf3Spec).



## Docker Image

[Docker Image for Iperf3](https://hub.docker.com/r/xridge/iperf3) is provided via [xridge's iperf3-docker repository](https://github.com/xridge/iperf3-docker).



## Legal

Iperf3 is licensed as 3-Clause BSD. 