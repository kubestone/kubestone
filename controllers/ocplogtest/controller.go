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

package ocplogtest

import (
	"context"
	"github.com/xridge/kubestone/pkg/k8s"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Reconciler reconciles a OcpLogtest object
type Reconciler struct {
	K8S k8s.Access
	Log logr.Logger
}

// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=ocplogtests,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=ocplogtests/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=ocplogtests/finalizers,verbs=update

func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	var cr perfv1alpha1.OcpLogtest
	if err := r.K8S.Client.Get(ctx, req.NamespacedName, &cr); err != nil {
		return ctrl.Result{}, k8s.IgnoreNotFound(err)
	}

	if cr.Status.Completed {
		return ctrl.Result{}, nil
	}

	cr.Status.Running = true
	if err := r.K8S.Client.Status().Update(ctx, &cr); err != nil {
		return ctrl.Result{}, err
	}

	job := NewJob(&cr)
	if err := r.K8S.CreateWithReference(ctx, job, &cr); err != nil {
		return ctrl.Result{}, err
	}

	jobFinished, err := r.K8S.IsJobFinished(types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Name,
	})

	if err != nil {
		return ctrl.Result{}, err
	}
	if !jobFinished {
		// Wait for the job to be completed
		return ctrl.Result{Requeue: true}, nil
	}

	// The cr could have been modified since the last time we got it
	if err := r.K8S.Client.Get(ctx, req.NamespacedName, &cr); err != nil {
		return ctrl.Result{}, k8s.IgnoreNotFound(err)
	}
	cr.Status.Running = false
	cr.Status.Completed = true
	if err := r.K8S.Client.Status().Update(ctx, &cr); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&perfv1alpha1.OcpLogtest{}).
		Complete(r)
}
