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

// Iperf3ConfigurationSpec contains configuration parameters
// with scheduling options for the both the iperf3 client
// and server instances.
type Iperf3ConfigurationSpec struct {
	// Command line arguments appended to the
	// predefined iperf3 parameters
	// +optional
	CmdLineArgs string `json:"cmdLineArgs"`

	// Labels to add to the iperf3 pod.
	// +optional
	PodLabels map[string]string `json:"podLabels,omitempty"`

	// Scheduling related options to determine which
	// node the iperf3 pod should be scheduled
	// +optional
	PodScheduling PodSchedulingSpec `json:"podScheduling,omitempty"`

	// Host Network requested for the iperf3 pod, if enabled the
	// hosts network namespace is used. Default to false.
	// +optional
	HostNetwork bool `json:"hostNetwork,omitempty"`
}

// Iperf3Spec defines the Iperf3 Benchmark Stone which
// consist of server deployment with service definition
// and client pod.
type Iperf3Spec struct {
	// Image defines the iperf3 docker image used for the benchmark
	Image ImageSpec `json:"image"`

	// ServerConfiguration contains the configuration of the iperf3 server
	// +optional
	ServerConfiguration Iperf3ConfigurationSpec `json:"serverConfiguration,omitempty"`

	// ClientConfiguration contains the configuration of the iperf3 client
	// +optional
	ClientConfiguration Iperf3ConfigurationSpec `json:"clientConfiguration,omitempty"`

	// Use UDP rather than TCP. If enabled the '--udp' parameter is set.
	// +optional
	UDP bool `json:"udp,omitempty"`
}

// Iperf3Status describes the current state of the benchmark
type Iperf3Status struct {
	// State of execution
	Running bool `json:"running"`
	// State of completion
	Completed bool `json:"completed"`
}

// +kubebuilder:object:root=true

// Iperf3 is the Schema for the iperf3s API
type Iperf3 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Iperf3Spec   `json:"spec,omitempty"`
	Status Iperf3Status `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// Iperf3List contains a list of Iperf3
type Iperf3List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Iperf3 `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Iperf3{}, &Iperf3List{})
}
