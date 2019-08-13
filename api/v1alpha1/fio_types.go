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

// FioSpec defines the desired state of Fio
type FioSpec struct {
	// Image defines the fio docker image used for the benchmark
	Image ImageSpec `json:"image"`

	// BuiltinJobFiles contains a list of fio job files that are already present
	// in the docker image
	// +optional
	BuiltinJobFiles []string `json:"builtinJobFiles,omitempty"`

	// TODO: Add implementation for custom job files (job as string in CR)

	// CmdLineArgs are appended to the predefined fio parameters
	// +optional
	CmdLineArgs string `json:"cmdLineArgs,omitempty"`
}

// FioStatus describes the current state of the benchmark
type FioStatus struct {
	// Running shows the state of execution
	Running bool `json:"running"`
	// Completed shows the state of completion
	Completed bool `json:"completed"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Running",type="boolean",JSONPath=".status.running"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"

// Fio is the Schema for the fios API
type Fio struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FioSpec   `json:"spec,omitempty"`
	Status FioStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FioList contains a list of Fio
type FioList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Fio `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Fio{}, &FioList{})
}
