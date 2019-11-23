title: Kubestone - Sysbench: Database and System performance benchmark

# Sysbench - Scriptable database and system performance benchmark 

!!! quote
    sysbench is a scriptable multi-threaded benchmark tool based on LuaJIT. It is most frequently used for database benchmarks, but can also be used to create arbitrarily complex workloads that do not involve a database server. 

With the [sysbench](https://github.com/akopytov/sysbench) benchmark you can measure the CPU, Memory, Database and Filesystem characteritics of your Kubernetes cluster. 



## Mode of operation

Kubestone generates a Kubernetes Job from each Sysbench CR that will run a single pod with the defined job. Sysbench's input parameters can be specified in the CR with their respective names:
`sysbench [options]... [testname] [command]`



## Example configuration

You can find [configuration example](https://github.com/xridge/kubestone/blob/master/config/samples/perf_v1alpha1_sysbench.yaml) in the GitHub repository.




## Sample benchmark
```bash
kubectl create --namespace kubestone -f https://raw.githubusercontent.com/xridge/kubestone/master/config/samples/perf_v1alpha1_sysbench.yaml
```


Please refer to the [quickstart guide](../quickstart.md) for details on generic principles and setup of Kubestone.




## Sysbench Configuration

The complete documentation of sysbench CR can be found in the [API Docs](../apidocs.md#perf.kubestone.xridge.io/v1alpha1.SysbenchSpec).



## Docker Image

[Docker Image for sysbench](https://hub.docker.com/r/xridge/sysbench) is provided via [xridge's sysbench repository](https://github.com/xridge/sysbench-docker).



## Legal

sysbench is licensed as GPLv2. 