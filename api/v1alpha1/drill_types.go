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

// DrillSpec defines a benchmark run for drill load tester
type DrillSpec struct {
	// Image defines the drill docker image used for the benchmark
	Image ImageSpec `json:"image"`

	// BenchmarksVolume contains the files describing the benchmarks.
	// A ConfigMap is created from this map, where the keys will be used
	// as filenames and values will represent the contents of the files.
	BenchmarksVolume map[string]string `json:"benchmarksVolume"`

	// BenchmarkFile is the top level file (entry point) specified to drill.
	BenchmarkFile string `json:"benchmarkFile"`

	// Options are appended to the predefined drill parameters
	// +optional
	Options string `json:"options,omitempty"`

	// PodConfig contains the configuration for the benchmark pod, including
	// pod labels and scheduling policies (affinity, toleration, node selector...)
	// +optional
	PodConfig PodConfigurationSpec `json:"podConfig,omitempty"`
}

// DrillStatus describes the current state of the benchmark
type DrillStatus struct {
	// Running shows the state of execution
	Running bool `json:"running"`
	// Completed shows the state of completion
	Completed bool `json:"completed"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Running",type="boolean",JSONPath=".status.running"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"

// Drill is the Schema for the drills API
type Drill struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DrillSpec   `json:"spec,omitempty"`
	Status DrillStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DrillList contains a list of Drill
type DrillList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Drill `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Drill{}, &DrillList{})
}
