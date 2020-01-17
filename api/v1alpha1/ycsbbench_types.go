package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// YcsbBenchSpec defines the desired state of YcsbBench
type YcsbBenchSpec struct {
	// Image defines the docker image used for the benchmark
	Image ImageSpec `json:"image"`

	Action     string `json:"action"`
	DbType     string `json:"db_type"`
	Workletter string `json:"workletter"`
	DbArgs     string `json:"db_args"`

	// PodConfig contains the configuration for the benchmark pod, including
	// pod labels and scheduling policies (affinity, toleration, node selector...)
	// +optional
	PodConfig PodConfigurationSpec `json:"podConfig,omitempty"`
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
