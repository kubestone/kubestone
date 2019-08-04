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

package iperf3

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"

	corev1 "k8s.io/api/core/v1"

	"github.com/xridge/kubestone/pkg/k8s"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// Reconciler provides fields from manager to reconciler
type Reconciler struct {
	K8S k8s.Access
	Log logr.Logger
}

// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=iperf3s,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=iperf3s/status,verbs=get;update;patch

// Reconcile Iperf3 Job Requests
func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("iperf3", req.NamespacedName)

	var cr perfv1alpha1.Iperf3
	if err := r.K8S.Client.Get(ctx, req.NamespacedName, &cr); err != nil {
		//log.Error(err, "Unable to fetch Iperf3 CR")
		return ctrl.Result{}, k8s.IgnoreNotFound(err)
	}

	crRef, err := reference.GetReference(r.K8S.Scheme, &cr)
	if err != nil {
		log.Error(err, "Unable to get reference to Iperf3 CR")
		return ctrl.Result{}, err
	}

	if err := r.newServerDeployment(ctx, &cr, crRef); err != nil {
		r.K8S.EventRecorder.Eventf(crRef, corev1.EventTypeWarning, k8s.CreateFailed,
			"Unable to create iperf3 Server Deployment: %v", err)
		return ctrl.Result{}, err
	}

	if err := r.newServerService(ctx, &cr, crRef); err != nil {
		r.K8S.EventRecorder.Eventf(crRef, corev1.EventTypeWarning, k8s.CreateFailed,
			"Unable to create iperf3 Server Service: %v", err)
		return ctrl.Result{}, err
	}

	serverReady, err := r.serverDeploymentReady(&cr)
	if err != nil {
		r.K8S.EventRecorder.Eventf(crRef, corev1.EventTypeWarning, k8s.CreateFailed,
			"Unable to determine iperf3 Server Deployment state: %v", err)
		return ctrl.Result{}, err
	}
	if !serverReady {
		// Wait for the client pod to the deployment to be ready
		return ctrl.Result{Requeue: true}, nil
	}

	if err := r.newClientPod(ctx, &cr, crRef); err != nil {
		r.K8S.EventRecorder.Eventf(crRef, corev1.EventTypeWarning, k8s.CreateFailed,
			"Unable to create iperf3 client pod: %v", err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager registers the Iperf3Reconciler with the provided manager
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&perfv1alpha1.Iperf3{}).
		Complete(r)
}
