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

// YcsbBenchSpec defines the desired state of YcsbBench
type YcsbBenchSpec struct {
	// Image defines the docker image used for the benchmark
	Image ImageSpec `json:"image"`

	Database string `json:"database"`
	Workload string `json:"workload"`
	// +optional
	Options    YcsbBenchOptions  `json:"options,omitempty"`
	Properties map[string]string `json:"properties"`

	// PodConfig contains the configuration for the benchmark pod, including
	// pod labels and scheduling policies (affinity, toleration, node selector...)
	// +optional
	PodConfig PodConfigurationSpec `json:"podConfig,omitempty"`
}

type YcsbBenchOptions struct {
	Threadcount int `json:"threadcount,omitempty"`
	Target      int `json:"target,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Running",type="boolean",JSONPath=".status.running"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"

// YcsbBench is the Schema for the ycsbbenches API
type YcsbBench struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   YcsbBenchSpec   `json:"spec,omitempty"`
	Status BenchmarkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// YcsbBenchList contains a list of YcsbBench
type YcsbBenchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []YcsbBench `json:"items"`
}

func init() {
	SchemeBuilder.Register(&YcsbBench{}, &YcsbBenchList{})
}
