![Kubestone](images/kubestone-logo-notext.png)  
# Kubestone

Kubestone is a benchmarking [Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) that can evaluate the performance of [Kubernetes](https://kubernetes.io) installations. 



## Features

- **Supports common set of benchmarks** to measure:
  CPU, Disk, Network and Application performance
- **Fine-grained control over Kubernetes Scheduling primitives**:
  Affinity, Anti-Affinity, Tolerations, Storage Classes and Node Selection  
- **Cloud Native**: 
  Benchmarks runs are defined as Custom Resources and executed in the cluster using  Kubernetes resources: Pods, Jobs, Deployments and Services.
- **Extensible**: 
  New benchmarks can easily be added by implementing a new controller. 



## Benchmarks

{!benchmark-status.md!}


Follow the [quickstart guide](quickstart.md) to see how Kubestone can be deployed and how benchmarks can be run.



## Community

You can reach us on [Slack](https://join.slack.com/t/kubestone/signup) and via the [Kubestone Mail Group](https://groups.google.com/forum/#!forum/kubestone). 



## Contributing

All contributions are welcome! Bug reports, fixes, new features, documentation improvements and ideas help us to create the most comprehensive benchmark suite for Kubernetes. 

If you would like to get involved please read the [development guide](devguide.md). 

Issues labelled with '[good first issue](https://github.com/xridge/kubestone/labels/good%20first%20issue)' and '[help wanted](https://github.com/xridge/kubestone/labels/help%20wanted)' in [Kubestone repository](https://github.com/xridge/kubestone) are good starting points to join the community.



For long term plans please refer to the [Projects](https://github.com/xridge/kubestone/projects) and [Milestones](https://github.com/xridge/kubestone/milestones) pages.



## License

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at 

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.