# Iperf3 - Network bandwidth benchmark

!!! quote
   iPerf3 is a tool for active measurements of the maximum achievable bandwidth on IP networks. It supports tuning of various parameters related to timing, buffers and protocols (TCP, UDP, SCTP with IPv4 and IPv6). 

With the [iperf3](https://iperf.fr/) benchmark you can measure the I/O performance of the network hardware and stack used in your Kubernetes cluster. 



## Mode of operation

TBD

## Sample benchmark
TBD

Please refer to the [quickstart guide](quickstart.md) for details on generic principles and setup of Kubestone.




## Reference configuration
You can find [configuration example](https://github.com/xridge/kubestone/blob/master/config/samples/perf_v1alpha1_iperf3.yaml) in the GitHub repository.



## Further documentation

The complete documentation of iperf3 (and other benchmarks') CR can be found in the API Docs.



## Docker Image

[Docker Image for Iperf3](https://hub.docker.com/r/xridge/iperf3) is provided via [xridge's iperf3-docker repository](https://github.com/xridge/iperf3-docker).



## Legal

Iperf3 is licensed as 3-Clause BSD. 