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
	// JobFiles lists the file names of the job files that fio should run
	// +kubebuilder:validation:MinItems=1
	JobFiles []string `json:"jobFiles"`

	// RemoteJobFiles lists the URLs of optional remote job files. All the given
	// files will be downloaded and you can use the `jobFiles` field to select which
	// ones to run
	RemoteJobFiles []string `json:"remoteJobFiles,omitempty"`
}

// FioStatus defines the observed state of Fio
type FioStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

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
