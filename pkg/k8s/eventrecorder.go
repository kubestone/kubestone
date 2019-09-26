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

package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"

	corev1 "k8s.io/api/core/v1"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	// CreateFailed is an event provided via EventRecorder
	CreateFailed = "CreateFailed"
	// Created is an event provided via EventRecorder
	Created = "Created"
	// Deleted is an event provided via EventRecorder
	Deleted = "Deleted"
)

// NewEventRecorder creates a new event recorder
func NewEventRecorder(clientSet *kubernetes.Clientset, logf func(format string, args ...interface{})) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	if logf != nil {
		eventBroadcaster.StartLogging(logf)
	}
	eventBroadcaster.StartRecordingToSink(
		&typedcorev1.EventSinkImpl{Interface: clientSet.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(k8sscheme.Scheme,
		corev1.EventSource{Component: "kubestone"})
	return recorder
}
