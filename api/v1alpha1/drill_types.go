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

// DrillSpec defines benchmark run for drill load tester
// The benchmarkFile, and options is passed to drill as follows:
// drill [OPTIONS] --benchmark <benchmarkFile>
type DrillSpec struct {
	// Image defines the drill docker image used for the benchmark
	Image ImageSpec `json:"image"`

	// BenchmarksVolume holds the content of benchmark files.
	// The key of the map specifies the filename and the value is the content
	// of the file. ConfigMap is created from the map which is mounted as
	// benchmarks directory to the benchmark pod.
	BenchmarksVolume map[string]string `json:"benchmarksVolume"`

	// BenchmarkFile is the entry point file (passed to --benchmark) specified to drill.
	BenchmarkFile string `json:"benchmarkFile"`

	// Options are appended to the options parameter set of drill
	// +optional
	Options string `json:"options,omitempty"`

	// PodConfig contains the configuration for the benchmark pod, including
	// pod labels and scheduling policies (affinity, toleration, node selector...)
	// +optional
	PodConfig PodConfigurationSpec `json:"podConfig,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Running",type="boolean",JSONPath=".status.running"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"

// Drill is the Schema for the drills API
type Drill struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DrillSpec       `json:"spec,omitempty"`
	Status BenchmarkStatus `json:"status,omitempty"`
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
