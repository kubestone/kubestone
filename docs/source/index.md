# Kubestone
Kubestone is a benchmarking [Kubernetes Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) that can evaluate the performance of Kubernetes Clusters. 



## Features

- **Supports common set of benchmarks** to measure:
  CPU, Disk and Network and Application performance

- **Fine-grained control on Kubernetes Scheduling primitives**:
  Affinity, Anti-Affinity, Tolerations, Storage Classes and Node Selection  

- **Cloud Native benchmarking**: 

  Runs are defined as Custom Resources. Benchmarks are executed within the cluster using Kubernetes resources: Pods, Jobs, Deployments and Services.

- **Extensible**: 
  New benchmarks can easily be added by implementing a new controller. 



## Benchmarks

| Type              | Benchmark name | Status                                                       |
| ----------------- | :------------: | ------------------------------------------------------------ |
| Core/CPU          |    Sysbench    | [Under development](https://github.com/xridge/kubestone/pull/71) |
| Core/Disk         |      Fio       | [Supported](https://github.com/xridge/kubestone/blob/master/config/samples/fio/base/fio_cr.yaml) |
| Core/Network      |     Iperf3     | [Supported](https://github.com/xridge/kubestone/blob/master/config/samples/perf_v1alpha1_iperf3.yaml) |
| Application/Etcd  |      Etcd      | [Planned](https://github.com/xridge/kubestone/issues/15)     |
| Application/K8S   |    KubePerf    | [Planned](https://github.com/xridge/kubestone/issues/14)     |
| Application/Spark |   SparkBench   | Under development                                            |



## Installation

### Requirements

- [Kubernetes](https://kubernetes.io/) v1.14 (or newer)
- [Kustomize](https://kustomize.io/) v3.10

### Steps





### Community

Slack channel comes here

Mail group comes here too



### Contributing

All contributions are welcome! Bug reports, fixes, new features, documentation improvements and ideas help us jointly create the most comprehensive benchmark suite for Kubernetes. 

If you would like to get involved please read through on the development guide. 

Issues labelled with '[good first issue](https://github.com/xridge/kubestone/labels/good%20first%20issue)' and '[help wanted](https://github.com/xridge/kubestone/labels/help%20wanted)' are good starting points to join the community.



### License

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at 
http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.