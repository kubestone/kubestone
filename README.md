![](https://raw.githubusercontent.com/xridge/kubestone/master/images/kubestone-logo.png)


[![CircleCI](https://circleci.com/gh/xridge/kubestone/tree/master.svg?style=shield)](https://circleci.com/gh/xridge/kubestone/tree/master)
[![docker build](https://img.shields.io/docker/cloud/build/xridge/kubestone.svg)](https://hub.docker.com/r/xridge/kubestone)
[![docker pulls](https://img.shields.io/docker/pulls/xridge/kubestone.svg)](https://hub.docker.com/r/xridge/kubestone)
[![Go Report Card](https://goreportcard.com/badge/github.com/xridge/kubestone)](https://goreportcard.com/report/github.com/xridge/kubestone)
[![license](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

#
kubestone is a benchmarking Kubernetes Operator. 

It provides cpu, memory and disk performance measurements
for Kubernetes cluster via common set of benchmarks:
 * Disk: [fio](https://fio.readthedocs.io)
 * Network: [IPerf3](https://iperf.fr)
 * CPU: [sysbench](https://wiki.gentoo.org/wiki/Sysbench)
 
## Benchmark definitions
Benchmark are initiated by creating Custom Resources in
any namespace in Kubernetes. When a new Kubestone CR is created
the benchmark's workflow is executed.

## Usage
### Prerequisities
 * [Kubernetes](https://kubernetes.io) v1.14 (or newer)
 * [Kustomzie](https://kustomize.io) v3.10
 * [Go](https://golang.org) v1.12 (with `GO111MODULE=on`)


### Install the Custom Resource Definitions
 * Install the CRDs to Kubernetes
    ```bash
     $ make install
        /Users/dev/goenvs/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
        kubectl apply -f config/crd/bases
        customresourcedefinition.apiextensions.k8s.io/fios.perf.kubestone.xridge.io configured
        customresourcedefinition.apiextensions.k8s.io/iperf3s.perf.kubestone.xridge.io configured
    ```
 * Deploy the Operator
    ```bash
     $ make deploy
        /Users/devs/goenvs/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
        kubectl apply -f config/crd/bases
        customresourcedefinition.apiextensions.k8s.io/fios.perf.kubestone.xridge.io created
        customresourcedefinition.apiextensions.k8s.io/iperf3s.perf.kubestone.xridge.io created
        kustomize build config/default | kubectl apply -f -
        namespace/kubestone-system created
        customresourcedefinition.apiextensions.k8s.io/fios.perf.kubestone.xridge.io configured
        customresourcedefinition.apiextensions.k8s.io/iperf3s.perf.kubestone.xridge.io configured
        role.rbac.authorization.k8s.io/kubestone-leader-election-role created
        clusterrole.rbac.authorization.k8s.io/kubestone-manager-role created
        clusterrole.rbac.authorization.k8s.io/kubestone-proxy-role created
        rolebinding.rbac.authorization.k8s.io/kubestone-leader-election-rolebinding created
        clusterrolebinding.rbac.authorization.k8s.io/kubestone-manager-rolebinding created
        clusterrolebinding.rbac.authorization.k8s.io/kubestone-proxy-rolebinding created
        service/kubestone-controller-manager-metrics-service created
        deployment.apps/kubestone-controller-manager created
    ```

### Run benchmark
 * Create dedicated namespace for bencmarking
 ```bash
  $ kubectl create namespace kubestone
    namespace/kubestone created
 ```
 * Start sample benchmark by creating CR
 ```bash
  $ kubectl create -n kubestone -f config/samples/perf_v1alpha1_iperf3.yaml
    iperf3.perf.kubestone.xridge.io/iperf3-sample created
 ```

Sample benchmarks are located in [config/samples/](config/samples).

### Inspect benchmark
Benchmarks are executed within the same namespace where the CR is created.
```bash
  $ kubectl get all -n kubestone
    NAME                                READY   STATUS    RESTARTS   AGE
    pod/iperf3-sample-6cb5445f7-mxhd2   1/1     Running   0          5s
    pod/iperf3-sample-client            1/1     Running   0          2s


    NAME                    TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
    service/iperf3-sample   ClusterIP   10.106.76.97   <none>        5201/TCP   5s


    NAME                            READY   UP-TO-DATE   AVAILABLE   AGE
    deployment.apps/iperf3-sample   1/1     1            1           5s

    NAME                                      DESIRED   CURRENT   READY   AGE
    replicaset.apps/iperf3-sample-6cb5445f7   1         1         1       5s

  $ kubctl logs -n kubestone iperf3-sample-client
    Connecting to host iperf3-sample, port 5201
    [  5] local 10.233.116.159 port 43028 connected to 10.233.60.64 port 5201
    [ ID] Interval           Transfer     Bitrate         Retr  Cwnd
    [  5]   0.00-1.00   sec   112 MBytes   940 Mbits/sec  159    492 KBytes
    [  5]   1.00-2.00   sec   109 MBytes   912 Mbits/sec  156    301 KBytes
    [  5]   2.00-3.00   sec   110 MBytes   923 Mbits/sec    0    422 KBytes
    [  5]   3.00-4.00   sec   109 MBytes   912 Mbits/sec    8    404 KBytes
    [  5]   4.00-5.00   sec   109 MBytes   912 Mbits/sec    0    571 KBytes
    [  5]   5.00-6.00   sec   110 MBytes   923 Mbits/sec   46    491 KBytes
    [  5]   6.00-7.00   sec   109 MBytes   912 Mbits/sec    0    636 KBytes
    [  5]   7.00-8.00   sec   109 MBytes   912 Mbits/sec   46    569 KBytes
    [  5]   8.00-9.00   sec   110 MBytes   923 Mbits/sec   51    339 KBytes
    [  5]   9.00-10.00  sec   110 MBytes   923 Mbits/sec    0    526 KBytes
    - - - - - - - - - - - - - - - - - - - - - - - - -
    [ ID] Interval           Transfer     Bitrate         Retr
    [  5]   0.00-10.00  sec  1.07 GBytes   919 Mbits/sec  466             sender
    [  5]   0.00-10.04  sec  1.07 GBytes   913 Mbits/sec                  receiver

    iperf Done.

```


## Development status
Under development / alpha.


## License
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
