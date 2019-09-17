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

// QperfPort is the TCP port where the qperf server and client listens
const QperfPort = 19765

// QperfConfigurationSpec contains configuration parameters
// with scheduling options for the both the qperf client
// and server instances.
type QperfConfigurationSpec struct {
	PodConfigurationSpec `json:",inline"`

	// HostNetwork requested for the qperf pod, if enabled the
	// hosts network namespace is used. Default to false.
	// +optional
	HostNetwork bool `json:"hostNetwork,omitempty"`
}

// QperfSpec defines the Qperf Benchmark Stone which
// consist of server deployment with service definition
// and client pod.
type QperfSpec struct {
	// Image defines the qperf docker image used for the benchmark
	Image ImageSpec `json:"image"`

	// Options are options for the qperf binary
	// +optional
	Options string `json:"options,omitempty"`

	// Tests are the tests that we would like to run
	Tests []string `json:"tests"`

	// ServerConfiguration contains the configuration of the qperf server
	// +optional
	ServerConfiguration QperfConfigurationSpec `json:"serverConfiguration,omitempty"`

	// ClientConfiguration contains the configuration of the qperf client
	// +optional
	ClientConfiguration QperfConfigurationSpec `json:"clientConfiguration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Running",type="boolean",JSONPath=".status.running"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"

// Qperf is the Schema for the qperves API
type Qperf struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QperfSpec       `json:"spec,omitempty"`
	Status BenchmarkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// QperfList contains a list of Qperf
type QperfList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Qperf `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Qperf{}, &QperfList{})
}
