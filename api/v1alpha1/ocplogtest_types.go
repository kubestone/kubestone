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

// OcpLogtestSpec defines the desired state of OcpLogtest
type OcpLogtestSpec struct {
	// Image defines the docker image used for the benchmark
	Image ImageSpec `json:"image"`

	// length of each line
	LineLength int `json:"lineLength,omitempty"`

	// number of lines to generate
	NumLines int `json:"numLines,omitempty"`

	// lines per minute
	Rate int `json:"rate,omitempty"`

	// repeat the same line of text over and over or use new text for each line
	FixedLine bool `json:"fixedLine,omitempty"`

	// PodConfig contains the configuration for the benchmark pod, including
	// pod labels and scheduling policies (affinity, toleration, node selector...)
	// +optional
	PodConfig PodConfigurationSpec `json:"podConfig,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Running",type="boolean",JSONPath=".status.running"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"

// OcpLogtest is the Schema for the ocplogtests API
type OcpLogtest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OcpLogtestSpec  `json:"spec,omitempty"`
	Status BenchmarkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OcpLogtestList contains a list of OcpLogtest
type OcpLogtestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OcpLogtest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OcpLogtest{}, &OcpLogtestList{})
}
