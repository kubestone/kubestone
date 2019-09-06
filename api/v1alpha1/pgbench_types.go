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

// PostgresSpec contains the configuration parameters for the PostreSQL database
type PostgresSpec struct {
	// Host is the name of host to connect to
	Host string `json:"host"`

	// Port number to connect to at the server host
	Port int `json:"port"`

	// User is the PostgreSQL user name to connect as
	User string `json:"user"`

	// Password is to be used if the server demands password authentication
	Password string `json:"password"`

	// Database is name of the database
	Database string `json:"database"`
}

// PgbenchSpec describes a pgbench benchmark job
type PgbenchSpec struct {
	// Image defines the docker image used for the benchmark
	Image ImageSpec `json:"image"`

	// Postgres contains the configuration parameters for the PostgreSQL database
	// that will run the benchmark
	Postgres PostgresSpec `json:"postgres"`

	// InitArgs contains the command line arguments passed to the init container
	// +optional
	InitArgs string `json:"initArgs,omitempty"`

	// Args contains the command line arguments passed to the main pgbench container
	// +optional
	Args string `json:"args,omitempty"`

	// PodConfig contains the configuration for the benchmark pod, including
	// pod labels and scheduling policies (affinity, toleration, node selector...)
	// +optional
	PodConfig PodConfigurationSpec `json:"podConfig,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Running",type="boolean",JSONPath=".status.running"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"

// Pgbench is the Schema for the pgbenches API
type Pgbench struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PgbenchSpec     `json:"spec,omitempty"`
	Status BenchmarkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PgbenchList contains a list of Pgbench
type PgbenchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pgbench `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pgbench{}, &PgbenchList{})
}
