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

package pgbench

import (
	"context"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// Reconciler reconciles a Pgbench object
type Reconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=pgbenches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=pgbenches/status,verbs=get;update;patch

// Reconcile creates pgbench job(s) based on the custom resource(s)
func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("pgbench", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager registers the Reconciler with the provided manager
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&perfv1alpha1.Pgbench{}).
		Complete(r)
}
