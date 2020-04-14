/*
Copyright 2019 The xridge kubestone contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KafkaBenchSpec defines the desired state of KafkaBench
type KafkaBenchSpec struct {
	// Image defines the kafka docker image used for the benchmark
	Image ImageSpec `json:"image"`

	// PodConfig contains the configuration for the benchmark pod, including
	// pod labels and scheduling policies (affinity, toleration, node selector...)
	// +optional
	PodConfig PodConfigurationSpec `json:"podConfig,omitempty"`

	KafkaClusterInfo `json:",inline"`

	// Tests defines the tests with which to create
	Tests []KafkaTestSpec `json:"tests"`
}

// ClusterInfo to be used by the benchmark for ZooKeeper and Kafka Brokers
type KafkaClusterInfo struct {
	// List of ZooKeeper instances we to connect to
	ZooKeepers []string `json:"zookeepers"`

	// List of Kafka Broker instances we to connect to
	Brokers []string `json:"brokers"`
}

// TestSpec defines the specifications for the kafka tests
type KafkaTestSpec struct {
	Name        string `json:"name"`
	Threads     int32  `json:"threads"`
	Replication int    `json:"replication"`
	Partitions  int    `json:"partitions"`
	RecordSize  int    `json:"recordSize"`
	Records     int    `json:"records"`

	// ConsumerSleep defines the time in seconds the consumer will sleep before attempting to consume messages. Only change if you are having issues with consuming messages. Default: 40
	// +optional
	ConsumerSleep *int32 `json:"consumerSleep"`

	// Timeout defines the consumer maximum allowed time in milliseconds between returned records. (default: 10000)
	// +optional
	Timeout *int `json:"timeout"`

	// These can be any official producer Kafka options: https://kafka.apache.org/documentation/#producerconfigs
	// +optional
	ExtraProducerOpts []string `json:"extraProducerOpts"`
	ConsumersOnly     bool     `json:"consumersOnly"`
	ProducersOnly     bool     `json:"producersOnly"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Running",type="boolean",JSONPath=".status.running"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"

// KafkaBench is the Schema for the kafkabenches API
type KafkaBench struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KafkaBenchSpec  `json:"spec,omitempty"`
	Status BenchmarkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KafkaBenchList contains a list of KafkaBench
type KafkaBenchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KafkaBench `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KafkaBench{}, &KafkaBenchList{})
}
