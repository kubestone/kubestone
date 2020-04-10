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

// EsRallySpec defines the desired state of EsRally
type EsRallySpec struct {
	// Image defines the docker image used for the benchmark
	// +optional
	Image ImageSpec `json:"image,omitempty"`

	// PodConfig contains the configuration for the benchmark pod, including
	// pod labels and scheduling policies (affinity, toleration, node selector...)
	// +optional
	PodConfig PodConfigurationSpec `json:"podConfig,omitempty"`

	Track     string `json:"track"`
	Hosts     string `json:"hosts"`
	Pipeline  string `json:"pipeline"`
	Challenge string `json:"challenge"`

	Config      EsRallyConfig    `json:"config,omitempty"`
	Persistence EsRallyVolConfig `json:"persistence"`
}

type EsRallyVolConfig struct {
	Size         string `json:"size"`
	StorageClass string `json:"storageClass"`
}

type EsRallyConfig struct {
	// Nodes contains the number of esrally clients to use. Default is 1
	Nodes *int32 `json:"nodes"`
}

type EsRallyStatus struct {
	// Running shows the state of execution
	Running bool `json:"running"`
	// Completed shows the state of completion
	Completed bool `json:"completed"`
	// Deployed shows the state of the StatefulSet needed for testing
	Deployed bool `json:"deployed"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Deployed",type="boolean",JSONPath=".status.deployed"
// +kubebuilder:printcolumn:name="Running",type="boolean",JSONPath=".status.running"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"

// EsRally is the Schema for the esrallies API
type EsRally struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EsRallySpec   `json:"spec,omitempty"`
	Status EsRallyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EsRallyList contains a list of EsRally
type EsRallyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EsRally `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EsRally{}, &EsRallyList{})
}
