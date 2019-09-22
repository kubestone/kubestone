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
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Access provides client related structs to access kubernetes
type Access struct {
	Client        client.Client
	Clientset     *k8sclient.Clientset
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder
}

// RecordEventf is a convenience function to create an event (via Access.EventRecorder)
// for the given object.
func (a *Access) RecordEventf(object metav1.Object, eventtype, reason, messageFmt string, args ...interface{}) error {
	runtimeObject, ok := object.(runtime.Object)
	if !ok {
		return fmt.Errorf("object (%T) is not a runtime.Object", object)
	}

	objectRef, err := reference.GetReference(a.Scheme, runtimeObject)
	if err != nil {
		return fmt.Errorf("Unable to get reference to owner")
	}

	a.EventRecorder.Eventf(objectRef, eventtype, reason, messageFmt, args...)

	return nil
}

// CreateWithReference method creates a kubernetes resource and
// sets the owner reference to a given object. It provides basic
// idempotency (by ignoring Already Exists errors).
// Successful creation of the event is logged via EventRecorder
// to the owner.
func (a *Access) CreateWithReference(ctx context.Context, object, owner metav1.Object) error {
	runtimeObject, ok := object.(runtime.Object)
	if !ok {
		return fmt.Errorf("object (%T) is not a runtime.Object", object)
	}

	if err := controllerutil.SetControllerReference(owner, object, a.Scheme); err != nil {
		return err
	}

	err := a.Client.Create(ctx, runtimeObject)
	if IgnoreAlreadyExists(err) != nil {
		return err
	}

	if !errors.IsAlreadyExists(err) {
		_ = a.RecordEventf(owner, corev1.EventTypeNormal, Created,
			"Created %v", object.GetSelfLink())
	}

	return nil
}

// DeleteObject method deletes a kubernetes resource while
// ignores not found errors, so that it can be called multiple times.
// Successful deletion of the event is logged via EventRecorder
// to the owner.
func (a *Access) DeleteObject(ctx context.Context, object, owner metav1.Object) error {
	runtimeObject, ok := object.(runtime.Object)
	if !ok {
		return fmt.Errorf("object (%T) is not a runtime.Object", object)
	}

	// Need to get the object first so that the object.GetSelfLink()
	// works during Event Recording
	namespacedName := types.NamespacedName{
		Namespace: object.GetNamespace(),
		Name:      object.GetName(),
	}
	err := a.Client.Get(ctx, namespacedName, runtimeObject)
	if IgnoreNotFound(err) != nil {
		return err
	} else if errors.IsNotFound(err) {
		return nil
	}

	err = a.Client.Delete(ctx, runtimeObject)
	if IgnoreNotFound(err) != nil {
		return err
	}

	if !errors.IsNotFound(err) {
		_ = a.RecordEventf(owner, corev1.EventTypeNormal, Deleted,
			"Deleted %v", object.GetSelfLink())
	}

	return nil
}

// IsJobFinished returns true if the given job has already succeeded or failed
func (a *Access) IsJobFinished(namespacedName types.NamespacedName) (finished bool, err error) {
	job, err := a.Clientset.BatchV1().Jobs(namespacedName.Namespace).Get(
		namespacedName.Name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	finished = job.Status.CompletionTime != nil
	return finished, nil
}

// +kubebuilder:rbac:groups="",resources=endpoints,verbs=get;list

// IsEndpointReady returns true if the given endpoint is fully connected to at least one pod
func (a *Access) IsEndpointReady(namespacedName types.NamespacedName) (finished bool, err error) {
	// The Endpoint connection between the Service and the Pod is the final step before
	// a service becomes reachable in Kubernetes. When the endpoint is bound, your
	// service becomes connectable on vanilla k8s and azure, but not on GKE.
	// For details see #96: https://github.com/xridge/kubestone/issues/96
	//
	// Even though it is not enough to wait for the endpoints in certain cloud providers,
	// it is still the closest we can get between service creation and connectibility.
	endpoint, err := a.Clientset.CoreV1().Endpoints(namespacedName.Namespace).Get(
		namespacedName.Name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	readyAddresses := 0
	for _, subset := range endpoint.Subsets {
		if len(subset.NotReadyAddresses) > 0 {
			return false, nil
		}
		readyAddresses += len(subset.Addresses)
	}

	ready := readyAddresses > 0

	return ready, nil
}

// IsDeploymentReady returns true if the given deployment's ready replicas matching with the desired replicas
func (a *Access) IsDeploymentReady(namespacedName types.NamespacedName) (ready bool, err error) {
	ready, err = false, nil
	deployment, err := a.Clientset.AppsV1().Deployments(namespacedName.Namespace).Get(
		namespacedName.Name, metav1.GetOptions{})
	if err != nil {
		return ready, err
	}

	ready = deployment.Status.ReadyReplicas == *deployment.Spec.Replicas

	return ready, err
}
