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

// CreateWithReference method creates a kubernetes resource and
// sets the owner reference to a given object. It provides basic
// idempotency (by ignoring Already Exists errors).
// Successful creation of the event is logged via EventRecorder
// to the object.
func (a *Access) CreateWithReference(ctx context.Context, object, owner metav1.Object) error {
	runtimeObject, ok := object.(runtime.Object)
	if !ok {
		return fmt.Errorf("object (%T) is not a runtime.Object", object)
	}

	runtimeOwner, ok := owner.(runtime.Object)
	if !ok {
		return fmt.Errorf("owner (%T) is not a runtime.Object", object)
	}

	ownerRef, err := reference.GetReference(a.Scheme, runtimeOwner)
	if err != nil {
		return fmt.Errorf("Unable to get reference to owner")
	}

	if err := controllerutil.SetControllerReference(owner, object, a.Scheme); err != nil {
		return err
	}

	err = a.Client.Create(ctx, runtimeObject)
	if IgnoreAlreadyExists(err) != nil {
		return err
	}

	if !errors.IsAlreadyExists(err) {
		a.EventRecorder.Eventf(ownerRef, corev1.EventTypeNormal, CreateSucceeded,
			"Created %v", object.GetSelfLink())
	}

	return nil
}
