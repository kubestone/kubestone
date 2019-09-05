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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IopingVolumeSpec contains the configuration for the volume that the ioping job
// will use for benchmarking
// TODO: Factor VolumeSpec into a common one
type IopingVolumeSpec struct {
	// VolumeSource represents the source of the volume, e.g. an existing
	// PVC, host path, git repo, etc.
	VolumeSource corev1.VolumeSource `json:"volumeSource,omitempty"`

	// PersistentVolumeClaim describes the persistent volume claim that will be
	// created and used by the pod. This field *overrides* the VolumeSource to
	// point to the created PVC
	// +optional
	PersistentVolumeClaim *PersistentVolumeClaimSpec `json:"persistentVolumeClaim,omitempty"`
}

// IopingSpec defines the desired state of Ioping
type IopingSpec struct {
	// Image defines the ioping docker image used for the benchmark
	Image ImageSpec `json:"image"`

	// CmdLineArgs are appended to the predefined ioping parameters
	// +optional
	CmdLineArgs string `json:"cmdLineArgs,omitempty"`

	// PodConfig contains the configuration for the benchmark pod, including
	// pod labels and scheduling policies (affinity, toleration, node selector...)
	// +optional
	PodConfig PodConfigurationSpec `json:"podConfig,omitempty"`

	// Volume contains the configuration for the volume that the ioping job should
	// run on. If missing, no volume will attached to the job and Docker's layered
	// fs performance will be measured
	// +optional
	Volume *IopingVolumeSpec `json:"volume,omitempty"`
}

// IopingStatus describes the current state of the benchmark
type IopingStatus struct {
	// Running shows the state of execution
	Running bool `json:"running"`
	// Completed shows the state of completion
	Completed bool `json:"completed"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Running",type="boolean",JSONPath=".status.running"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"

// Ioping is the Schema for the iopings API
type Ioping struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IopingSpec   `json:"spec,omitempty"`
	Status IopingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IopingList contains a list of Ioping
type IopingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ioping `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Ioping{}, &IopingList{})
}
