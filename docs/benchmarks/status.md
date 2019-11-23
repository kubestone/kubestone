title: Kubestone - List of performance benchmarks

[//]: #(The links from this file are based on / as both files including this status page are located in /)

Current state of the benchmarks

| Type                    |           Benchmark name           | Status                                                                 |
| ----------------------- | :--------------------------------: | ---------------------------------------------------------------------- |
| Core/CPU                | [sysbench](benchmarks/sysbench.md) | [Supported](apidocs.md#perf.kubestone.xridge.io/v1alpha1.SysbenchSpec) |
| Core/Disk               |      [fio](benchmarks/fio.md)      | [Supported](apidocs.md#perf.kubestone.xridge.io/v1alpha1.FioSpec)      |
| Core/Disk               |   [ioping](benchmarks/ioping.md)   | [Supported](apidocs.md#perf.kubestone.xridge.io/v1alpha1.IopingSpec)   |
| Core/Memory             | [sysbench](benchmarks/sysbench.md) | [Supported](apidocs.md#perf.kubestone.xridge.io/v1alpha1.SysbenchSpec) |
| Core/Network            |   [iperf3](benchmarks/iperf3.md)   | [Supported](apidocs.md#perf.kubestone.xridge.io/v1alpha1.Iperf3Spec)   |
| Core/Network            |    [qperf](benchmarks/qperf.md)    | [Supported](apidocs.md#perf.kubestone.xridge.io/v1alpha1.QperfSpec)    |
| HTTP Load Tester        |    [drill](benchmarks/drill.md)    | [Supported](apidocs.md#perf.kubestone.xridge.io/v1alpha1.DrillSpec)    |
| Application/Etcd        |                etcd                | [Planned](https://github.com/xridge/kubestone/issues/15)               |
| Application/K8S         |              kubeperf              | [Planned](https://github.com/xridge/kubestone/issues/14)               |
| Application/PostgreSQL  |  [pgbench](benchmarks/pgbench.md)  | [Supported](apidocs.md#perf.kubestone.xridge.io/v1alpha1.PgbenchSpec)  |
| Application/Spark       |             sparkbench             | [Planned](https://github.com/xridge/kubestone/issues/83)               |
