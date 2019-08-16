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

	// PersistentVolumeClaimName is the name of the existing persistence volume
	// claim that will be used by the benchmark pod. If undefined, you can either
	// use PersistentVolumeClaim to create and use a PVC, or nothing to run the
	// benchmark inside the pod.
	// +optional
	PersistentVolumeClaimName *string `json:"persistentVolumeClaimName,omitempty"`

	// PersistentVolumeClaim describes the persistent volume claim that will be
	// created and used by the pod. This field is ignored if PersistentVolumeClaimName
	// is given, in that case the pod will use the PVC by that given name.
	// +optional
	PersistentVolumeClaim *PersistentVolumeClaimSpec `json:"persistentVolumeClaim,omitempty"`

	// BuiltinJobFiles contains a list of fio job files that are already present
	// in the docker image
	// +optional
	BuiltinJobFiles []string `json:"builtinJobFiles,omitempty"`

	// CustomJobFiles contains a list of custom fio job files
	// The exact format of fio job files is documented here:
	// https://fio.readthedocs.io/en/latest/fio_doc.html#job-file-format
	// The job files defined here will be mounted to the fio benchmark container
	// +optional
	CustomJobFiles []string `json:"customJobFiles,omitempty"`

	// CmdLineArgs are appended to the predefined fio parameters
	// +optional
	CmdLineArgs string `json:"cmdLineArgs,omitempty"`

	// TODO: Add affinity related structs (suggested by otto)
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
