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

// JMeterSpec defines the desired state of JMeter
type JMeterSpec struct {
	// Image defines the docker image used for the benchmark
	Image ImageSpec `json:"image"`

	// PlanTest define the jmeter plan test
	PlanTest map[string]string `json:"planTest"`

	// TestName define the jmeter test name
	TestName string `json:"testName"`

	// Properties files definitions
	// +optional
	Props map[string]string `json:"props,omitempty"`

	// Properties passed to jmeter
	// +optional
	PropsName string `json:"propsName,omitempty"`

	// Volume to mount at result path
	Volume VolumeSpec `json:"volume"`

	// Job configurations
	// +optional
	JobConfig JobConfig `json:"jobConfig,omitempty"`

	// Args are appended to the predefined jmeter parameters
	// +optional
	Args string `json:"args,omitempty"`

	// Command contains the command line passed to the main jmeter container
	// +optional
	Command string `json:"command,omitempty"`

	// pod labels and scheduling policies (affinity, toleration, node selector...)
	// +optional
	Configuration PodConfigurationSpec `json:"configuration,omitempty"`
}

type JobConfig struct {
	// Specifies the desired number of successfully finished pods the job should be run with.
	// Setting to nil means that the success of any pod signals the success of all pods, and
	// allows parallelism to have any positive value.
	// Setting to 1 means that parallelism is limited to 1 and the success of that pod signals
	// the success of the job.
	// +optional
	Completions *int32 `json:"completions,omitempty"`

	// Specifies the maximum desired number of pods the job should run at any given time.
	// The actual number of pods running in steady state will be less than this number when:
	// ((.spec.completions - .status.successful) < .spec.parallelism),
	// i.e. when the work left to do is less than max parallelism.
	// +optional
	Parallelism *int32 `json:"parallelism,omitempty"`
}

// JMeterStatus defines the observed state of JMeter
type JMeterStatus struct {
	// Running shows the state of execution
	Running bool `json:"running"`
	// Completed shows the state of completion
	Completed bool `json:"completed"`
	// Valid shows the state of the validation
	Valid bool `json:"valid"`
}

// +kubebuilder:object:root=true

// JMeter is the Schema for the jmeters API
type JMeter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JMeterSpec   `json:"spec,omitempty"`
	Status JMeterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JMeterList contains a list of JMeter
type JMeterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JMeter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JMeter{}, &JMeterList{})
}
