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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/xridge/kubestone/pkg/k8s"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// Iperf3Reconciler provides fields from manager to reconciler
type Iperf3Reconciler struct {
	K8S k8s.Access
	Log logr.Logger
}

// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=iperf3s,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=iperf3s/status,verbs=get;update;patch

// Reconcile Iperf3 Job Requests
func (r *Iperf3Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("iperf3", req.NamespacedName)

	var cr perfv1alpha1.Iperf3
	if err := r.K8S.Client.Get(ctx, req.NamespacedName, &cr); err != nil {
		log.Error(err, "Unable to fetch Iperf3 CR")
		return ctrl.Result{}, k8s.IgnoreNotFound(err)
	}

	crRef, err := reference.GetReference(r.K8S.Scheme, &cr)
	if err != nil {
		log.Error(err, "Unable to get reference to Iperf3 CR")
		return ctrl.Result{}, err
	}

	iperf3ServerJob := newIperf3ServerJob(&cr)
	if err := controllerutil.SetControllerReference(&cr, iperf3ServerJob, r.K8S.Scheme); err != nil {
		r.K8S.EventRecorder.Eventf(crRef, corev1.EventTypeWarning, k8s.CreateFailed,
			"Unable to create reference for Iperf3 Server Job: %v", err)
		return ctrl.Result{}, err
	}
	if err := r.K8S.Client.Create(ctx, iperf3ServerJob); err != nil {
		if k8s.IgnoreNotFound(err) != nil {
			r.K8S.EventRecorder.Eventf(crRef, corev1.EventTypeWarning, k8s.CreateFailed,
				"Error creating Iperf3 Server Job: %v", err)
			return ctrl.Result{}, err
		}
	}
	r.K8S.EventRecorder.Eventf(crRef, corev1.EventTypeNormal, k8s.CreateSucceeded,
		"Created Iperf3 Server Job. Name: %v, Namespace: %v", iperf3ServerJob.Name, iperf3ServerJob.Namespace)

	// TODO: Add Iperf3 Client job

	return ctrl.Result{}, nil
}

func newIperf3ServerJob(cr *perfv1alpha1.Iperf3) *batchv1.Job {
	// TODO: Implement the logic here
	return nil
}

// SetupWithManager registers the Iperf3Reconciler with the provided manager
func (r *Iperf3Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&perfv1alpha1.Iperf3{}).
		Complete(r)
}
