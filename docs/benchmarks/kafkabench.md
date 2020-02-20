title: Kubestone - Kafka Bench: Kafka consumer and producer performance benchmark


# Kafka Bench - Kafka benchmark

With the Kafka benchmark, you can do the following by passing messages through the cluster used in your Kubernetes cluster:

* Measuring read and/or write throughput.
* Stress testing the cluster based on specific parameters (such as message size).
* Load testing for the purpose of evaluating specific metrics or determining the impact of cluster configuration changes.
 
## Mode of operation

To get benchmarks for kafka, a producer and consumer is created to pass messages through the kafka cluster. To use this you must already have a kafka cluster deployed. 
Addresses for the brokers and zookeepers must be provided.

To effectively pass messages through the cluster for measurement the controller creates the following objects during the benchmark (for each test defined):

* Producer Job
* Consumer Job
 
At the first step, a producer job is created. It will create a topic and start queueing up messages. 
A consumer job is also created, however it the init job will sleep for 40 seconds then start consuming messages.

The jobs use the [kafka-*-perf-test](https://docs.cloudera.com/runtime/7.0.3/kafka-managing/topics/kafka-manage-cli-perf-test.html) tools provided in by Kafka to assist in benchmarking. 
  
In order to avoid measuring loopback performance, it is advised that you set the affinity and anti-affinity scheduling primitives for the benchmark. The provided sample benchmark shows how to avoid executing the client and the server on the same machine. For further documentation please refer to Kubernetes' [respective documentation](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/).

## Example configuration

You can find [configuration example](https://github.com/xridge/kubestone/blob/master/config/samples/perf_v1alpha1_kafkabench.yaml) in the GitHub repository.

## Sample benchmark

```yaml
apiVersion: perf.kubestone.xridge.io/v1alpha1
kind: KafkaBench
metadata:
  name: kafkabench-sample
spec:
  image:
    name: confluentinc/cp-kafka:5.2.1

  # List of ZooKeeper instances we want to connect to
  zookeepers:
    - kafka-demo-cp-zookeeper-0.kafka-demo-cp-zookeeper-headless.kafka-demo.svc.sol1.diamanti.com:2181
    - kafka-demo-cp-zookeeper-1.kafka-demo-cp-zookeeper-headless.kafka-demo.svc.sol1.diamanti.com:2181
    - kafka-demo-cp-zookeeper-2.kafka-demo-cp-zookeeper-headless.kafka-demo.svc.sol1.diamanti.com:2181

  # List of Kafka Broker instances we want to connect to
  brokers:
    - kafka-demo-cp-kafka-0.kafka-demo-cp-kafka-headless.kafka-demo.svc.sol1.diamanti.com:9092
    - kafka-demo-cp-kafka-1.kafka-demo-cp-kafka-headless.kafka-demo.svc.sol1.diamanti.com:9092
    - kafka-demo-cp-kafka-2.kafka-demo-cp-kafka-headless.kafka-demo.svc.sol1.diamanti.com:9092

  ## Define performance tests we want to run against the new cluster
  tests:
    - name: "noreplication"
      # This is the number of instances we will fire up of the kafka-producer/kafka-consumer binaries
      threads: 2
      replication: 1
      partitions: 16
      recordSize: 100
      records: 60000000
      # These can be any official producer Kafka options: https://kafka.apache.org/documentation/#producerconfigs
      extraProducerOpts:
        - "buffer.memory=671088640"
      consumersOnly: false
      producersOnly: false
      timeout: 10000
```

Please refer to the [quickstart guide](../quickstart.md) for details on generic principles and setup of Kubestone.

## IPerf3 Configuration

The complete documentation of iperf3 CR can be found in the [API Docs](../apidocs.md#perf.kubestone.xridge.io/v1alpha1.KafkaBenchSpec).
