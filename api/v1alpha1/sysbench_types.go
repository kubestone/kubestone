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

// SysbenchSpec contains the configuration parameters
// with scheduling options for the sysbench benchmark.
// The options, testName and command parameters are passed
// to the sysbench benchmarking application.
type SysbenchSpec struct {
	// Image defines the sysbench docker image used for the benchmark
	Image ImageSpec `json:"image"`

	// PodConfig contains the configuration for the benchmark pod, including
	// pod labels and scheduling policies (affinity, toleration, node selector...)
	// +optional
	PodConfig PodConfigurationSpec `json:"podConfig,inline"`

	// Options is a list of zero or more command line options starting with '--'.
	// +optional
	Options string `json:"options,omitempty"`

	// TestName is the name of a built-in test (e.g. `fileio`, `memory`, `cpu`, etc.), or a name of one of the bundled
	// Lua scripts (e.g. `oltp_read_only`), or a path to a custom Lua script.
	TestName string `json:"testName"`

	// Command is an optional argument that will be passed by sysbench to the built-in test or script specified with
	// TestName. Command defines the action that must be performed by the test. The list of available commands depends
	// on a particular test. Some tests also implement their own custom commands.
	// +optional
	Command string `json:"command,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Running",type="boolean",JSONPath=".status.running"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"

// Sysbench is the Schema for the sysbenches API
type Sysbench struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SysbenchSpec    `json:"spec,omitempty"`
	Status BenchmarkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SysbenchList contains a list of Sysbench
type SysbenchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Sysbench `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Sysbench{}, &SysbenchList{})
}
