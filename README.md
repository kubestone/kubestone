# kubestone
[![CircleCI](https://circleci.com/gh/xridge/kubestone/tree/master.svg?style=shield)](https://circleci.com/gh/xridge/kubestone/tree/master)
[![docker build](https://img.shields.io/docker/cloud/build/xridge/kubestone.svg)](https://hub.docker.com/r/xridge/kubestone)
[![docker pulls](https://img.shields.io/docker/pulls/xridge/kubestone.svg)](https://hub.docker.com/r/xridge/kubestone)
[![Go Report Card](https://goreportcard.com/badge/github.com/xridge/kubestone)](https://goreportcard.com/report/github.com/xridge/kubestone)
[![license](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

kubestone is a benchmark oriented Kubernetes Operator. 

It provides cpu, memory and disk performance measurements
within a Kubernetes cluster via common set of Linux Benchmarks:
 * Disk: FIO
 * Network: IPerf3
 * CPU: sysbench
 
 
## Development status
Under development / alpha.

## Usage
```bash
 $ make deploy
 $ kubectl create -f config/samples/perf_v1alpha1_iperf3.yaml
```


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
