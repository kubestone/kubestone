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
)

// PullPolicy controls how the docker images are downloaded
// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
// +kubebuilder:validation:Enum=Always;Never;IfNotPresent
type PullPolicy string

// ImageSpec defines parameters for docker image executed on Kubernetes
type ImageSpec struct {
	// Name is the Docker Image location including the tag
	Name string `json:"name"`

	// +optional
	PullPolicy PullPolicy `json:"pullPolicy,omitempty"`

	// PullSecret is an optional list of references to secrets
	// in the same namespace to use for pulling any of the images
	// +optional
	PullSecret string `json:"pullSecret,omitempty"`
}

// PodSchedulingSpec encapsulates the scheduling related
// fields of a Kubernetes Pod
type PodSchedulingSpec struct {
	// Affinity is a group of affinity scheduling rules.
	// +optional
	Affinity corev1.Affinity `json:"affinity,omitempty"`

	// If specified, the pod's tolerations.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// A node selector represents the union of the results of
	// one or more label queries over a set of nodes; that is,
	// it represents the OR of the selectors represented by the
	// node selector terms.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// NodeName is a request to schedule this pod onto a specific node. If it is non-empty,
	// the scheduler simply schedules this pod onto that node, assuming that it fits resource
	// requirements.
	// +optional
	NodeName string `json:"nodeName,omitempty"`
}

// PodConfigurationSpec contains the configuration for the benchmark pods
type PodConfigurationSpec struct {
	// PodLabels are added to the pod as labels.
	// +optional
	PodLabels map[string]string `json:"podLabels,omitempty"`

	// PodScheduling contains options to determine which
	// node the pod should be scheduled on
	// +optional
	PodScheduling PodSchedulingSpec `json:"podScheduling,omitempty"`
}
