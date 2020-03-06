title: Kubestone - Openshift Container Platform Log Test: Log generator

# ocp-logtest - Log generation for OCP benchmarking

!!! quote
    The ocp_logtest.py script is a flexible tool for creating pod logs in OpenShift. 
    
The [OCP Logtest](https://github.com/openshift/svt/blob/master/openshift_scalability/content/logtest/ocp_logtest-README.md) can log random or fixed text for any given line/word sizes and at any given rate. It can run forever, for a set number of messages or for a set period of time.

## Mode of operation

The OcpLogtest spec will define and launch a Job that starts the `ocp_logtest.py` script with the defined arguments.

## Example configuration

You can find [configuration example](https://github.com/xridge/kubestone/blob/master/config/samples/perf_v1alpha1_ocplogtest.yaml) in the GitHub repository.


## Sample benchmark
```bash
$ kubectl create --namespace kubestone -f https://raw.githubusercontent.com/xridge/kubestone/master/config/samples/perf_v1alpha1_ocplogtest.yaml
```

Please refer to the [quickstart guide](../quickstart.md) for details on generic principles and setup of Kubestone.


## OcpLogtest Configuration

The complete documentation of OcpLogtest CR can be found in the [API Docs](../apidocs.md#perf.kubestone.xridge.io/v1alpha1.Iperf3Spec).

For option definitions and functions please reference the OCP Logtest [README](https://github.com/openshift/svt/blob/master/openshift_scalability/content/logtest/ocp_logtest-README.md#complete-ocp_logtestpy-flags)


## Docker Image

[Docker Image for OcpLogtest](https://quay.io/repository/mffiedler/ocp-logtest?tag=latest&tab=tags) and maintained by [Openshift](https://github.com/openshift/svt)  

## Legal

The `ocp_logtest.py` script is licensed as Apache License 2.0. 