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

package esrally

import (
	"context"
	"github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/go-logr/logr"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// EsRallyReconciler reconciles a EsRally object
type Reconciler struct {
	K8S k8s.Access
	Log logr.Logger
}

// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;create;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;create;delete

// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=esrallies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=esrallies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=esrallies/finalizers,verbs=update

func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("esrally", req.NamespacedName)

	var cr perfv1alpha1.EsRally
	if err := r.K8S.Client.Get(ctx, req.NamespacedName, &cr); err != nil {
		return ctrl.Result{}, k8s.IgnoreNotFound(err)
	}

	// Return if its completed
	if cr.Status.Completed {
		return ctrl.Result{}, nil
	}

	namespaceName := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Name,
	}

	if cr.Spec.Image == (v1alpha1.ImageSpec{}) {
		cr.Spec.Image = v1alpha1.ImageSpec{
			Name:       "diamantisolutions/esrally:kubestone",
			PullPolicy: "Always",
		}
	}

	// If its not running, create job and mark it as running
	if !cr.Status.Running {
		return esRallyJobHandler(cr, r, ctx, namespaceName)
	}

	// Grab the job pod to pass to statefulset
	pods, err := r.K8S.GetJobPods(namespaceName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if len(pods.Items) == 0 || pods.Items[0].Status.PodIP == "" {
		logger.Info("waiting for pod ip")
		return ctrl.Result{Requeue: true}, nil
	}

	// Deploy statefulset
	if !cr.Status.Deployed {
		return esRallyDeployHandler(cr, r, ctx, pods.Items[0].Status.PodIP)
	}

	_, ready, _ := r.K8S.IsStatefulSetReady(namespaceName)
	if !ready {
		// We need to wait for the StatefulSet to be ready, so requeue
		return ctrl.Result{Requeue: true}, nil
	}

	jobFinished, err := r.K8S.IsJobFinished(namespaceName)
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

func esRallyDeployHandler(cr perfv1alpha1.EsRally, r *Reconciler, ctx context.Context, ip string) (ctrl.Result, error) {
	statefulSet, sError := NewStatefulSet(&cr, ip)
	if sError != nil {
		return ctrl.Result{}, sError
	}

	// Create service
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Port: 1900, TargetPort: intstr.FromString("transport")},
			},
			Selector: statefulSet.Spec.Selector.MatchLabels,
		},
		Status: corev1.ServiceStatus{},
	}

	// Create service
	if err := r.K8S.CreateWithReference(ctx, &service, &cr); err != nil {
		return ctrl.Result{}, err
	}

	// Create StatefulSet
	if err := r.K8S.CreateWithReference(ctx, statefulSet, &cr); err != nil {
		return ctrl.Result{}, err
	}

	// set Deployed as true
	cr.Status.Deployed = true
	if err := r.K8S.Client.Status().Update(ctx, &cr); err != nil {
		return ctrl.Result{}, err
	}

	// We need to wait for the StatefulSet to be ready, so requeue
	return ctrl.Result{Requeue: true}, nil
}

func esRallyJobHandler(cr perfv1alpha1.EsRally, r *Reconciler, ctx context.Context, namespaceName types.NamespacedName) (ctrl.Result, error) {
	job := NewJob(&cr)
	if err := r.K8S.CreateWithReference(ctx, job, &cr); err != nil {
		return ctrl.Result{}, err
	}

	// Create coordinator service
	//masterService := GetEsRallyCoordSvc(&cr, job.Spec.Template.Labels)
	//if err := r.K8S.CreateWithReference(ctx, &masterService, &cr); err != nil {
	//	return ctrl.Result{}, err
	//}

	// Mark it as running
	cr.Status.Running = true
	if err := r.K8S.Client.Status().Update(ctx, &cr); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

func GetEsRallyCoordSvc(cr *perfv1alpha1.EsRally, selectorLabels map[string]string) corev1.Service {
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-coordinator",
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "transport",
					Port:       1900,
					TargetPort: intstr.FromString("transport"),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Selector: selectorLabels,
		},
		Status: corev1.ServiceStatus{},
	}

	return service
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&perfv1alpha1.EsRally{}).
		Complete(r)
}
