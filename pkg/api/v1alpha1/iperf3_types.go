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

// Iperf3Spec defines the Iperf3 Benchmark Job
type Iperf3Spec struct {
	// time in seconds to transmit for
	Time int32 `json:"time,omitempty"`

	// Use UDP rather than TCP
	// +optional
	UDP bool `json:"udp,omitempty"`
}

// Iperf3Status defines the observed state of Iperf3
type Iperf3Status struct {
	// Shows if the benchmark is running
	Running bool `json:"running"`
	// Shows completion of benchmark
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
