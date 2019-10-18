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

package ioping

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	corev1 "k8s.io/api/core/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
)

// Reconciler provides fields from manager to reconciler
type Reconciler struct {
	K8S k8s.Access
	Log logr.Logger
}

// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=iopings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=iopings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=iopings/finalizers,verbs=update

// Reconcile creates ioping job based on the custom resource
func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	var cr perfv1alpha1.Ioping
	if err := r.K8S.Client.Get(ctx, req.NamespacedName, &cr); err != nil {
		return ctrl.Result{}, k8s.IgnoreNotFound(err)
	}

	// Run to one completion
	if cr.Status.Completed {
		return ctrl.Result{}, nil
	}

	// Validate on first entry
	if !cr.Status.Completed && !cr.Status.Running {
		if valid, err := IsCrValid(&cr); !valid {
			_ = r.K8S.RecordEventf(&cr, corev1.EventTypeWarning, k8s.CreateFailed,
				"CR validation failed: %v", err)

			// Do not requeue invalid CRs
			return ctrl.Result{}, nil
		}
	}

	cr.Status.Running = true
	if err := r.K8S.Client.Status().Update(ctx, &cr); err != nil {
		return ctrl.Result{}, err
	}

	if cr.Spec.Volume.PersistentVolumeClaimSpec != nil {
		pvc := k8s.NewPersistentVolumeClaim(*cr.Spec.Volume.PersistentVolumeClaimSpec,
			cr.Name, cr.Namespace)
		if err := r.K8S.CreateWithReference(ctx, pvc, &cr); err != nil {
			return ctrl.Result{}, err
		}
		// Change ClaimName (from GENERATED) to the PVC was created
		cr.Spec.Volume.VolumeSource.PersistentVolumeClaim.ClaimName = cr.Name
	}

	job := NewJob(&cr)
	if err := r.K8S.CreateWithReference(ctx, job, &cr); err != nil {
		return ctrl.Result{}, err
	}

	// Check if finished
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

// SetupWithManager registers the Reconciler with the provided manager
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&perfv1alpha1.Ioping{}).
		Complete(r)
}
