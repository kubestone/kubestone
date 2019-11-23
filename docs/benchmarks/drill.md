title: Kubestone - Drill: HTTP Load Tester

# Drill - HTTP Load Tester

!!! quote
    Drill is a HTTP load testing application written in Rust. The main goal for this project is to build a really lightweight tool as alternative to other that require JVM and other stuff. You can write brenchmark files, in YAML format, describing all the stuff you want to test. It was inspired by Ansible syntax because it is really easy to use and extend.


With the [drill](https://github.com/fcsonline/drill) load generator, you can create a load test plan and execute it against any Web Service inside or outside of your Kubernetes installation. 



## Mode of operation

Drill is executed as a Kubernete Job by Kubestone. The user provided benchmark files are stored in a ConfigMap. The top level benchmark file (specified via `benchmarkFile`) is used to start the execution.



## Example configuration

You can find [configuration example](https://github.com/xridge/kubestone/blob/master/config/samples/perf_v1alpha1_drill.yaml) in the GitHub repository.



## Sample benchmark
```bash
$ kubectl create --namespace kubestone -f https://raw.githubusercontent.com/xridge/kubestone/master/config/samples/perf_v1alpha1_drill.yaml
```


Please refer to the [quickstart guide](../quickstart.md) for details on generic principles and setup of Kubestone.




## Drill Configuration

The complete documentation of drill CR can be found in the [API Docs](../apidocs.md#perf.kubestone.xridge.io/v1alpha1.DrillSpec).



## Docker Image

[Docker Image for Drill](https://hub.docker.com/r/xridge/drill) is provided via [xridge's drill-docker repository](https://github.com/xridge/drill-docker).



## Legal

Drill is licensed as GPLv3. 
