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
	"strings"

	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
)

// FioReconciler provides fields from manager to reconciler
type FioReconciler struct {
	K8S k8s.Access
	Log logr.Logger
}

// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=fios,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=fios/status,verbs=get;update;patch

func (r *FioReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("fio", req.NamespacedName)

	var cr perfv1alpha1.Fio
	if err := r.K8S.Client.Get(ctx, req.NamespacedName, &cr); err != nil {
		return ctrl.Result{}, k8s.IgnoreNotFound(err)
	}

	crRef, err := reference.GetReference(r.K8S.Scheme, &cr)
	if err != nil {
		log.Error(err, "Unable to get reference to FIO CR")
		return ctrl.Result{}, err
	}

	if err := r.newJob(ctx, &cr, crRef); err != nil {
		r.K8S.EventRecorder.Eventf(crRef, corev1.EventTypeWarning, k8s.CreateFailed,
			"Unable to create FIO job: %v", err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *FioReconciler) newJob(ctx context.Context, cr *perfv1alpha1.Fio, crRef *corev1.ObjectReference) error {
	labels := map[string]string{
		"app":               "fio",
		"kubestone-cr-name": cr.Name,
	}

	env := []corev1.EnvVar{
		{Name: "JOB_FILES", Value: strings.Join(cr.Spec.JobFiles, " ")},
	}
	if len(cr.Spec.RemoteJobFiles) > 0 {
		env = append(env, corev1.EnvVar{
			Name: "REMOTE_JOB_FILES", Value: strings.Join(cr.Spec.RemoteJobFiles, " "),
		})
	}

	backoffLimit := int32(0)

	imagePullPolicy := corev1.PullPolicy(cr.Spec.Image.PullPolicy)
	if imagePullPolicy == "" {
		if strings.HasSuffix(cr.Spec.Image.Name, ":latest") {
			imagePullPolicy = corev1.PullAlways
		} else {
			imagePullPolicy = corev1.PullIfNotPresent
		}
	}
	imagePullSecrets := []corev1.LocalObjectReference{}
	if cr.Spec.Image.PullSecret != "" {
		imagePullSecrets = append(imagePullSecrets, corev1.LocalObjectReference{
			Name: cr.Spec.Image.PullSecret,
		})
	}

	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "fio",
							Image:           cr.Spec.Image.Name,
							ImagePullPolicy: imagePullPolicy,
							Env:             env,
						},
					},
					RestartPolicy:    corev1.RestartPolicyNever,
					ImagePullSecrets: imagePullSecrets,
				},
			},
			BackoffLimit: &backoffLimit,
		},
	}

	// FIXME: The next three statements are common between serverdeployment, serverservice and clientpod,
	// it would make sense to factor it into one function. For that a type should be found which applies
	// for both SetControllerReference and Create's object.
	if err := controllerutil.SetControllerReference(cr, &job, r.K8S.Scheme); err != nil {
		return err
	}
	if err := r.K8S.Client.Create(ctx, &job); k8s.IgnoreAlreadyExists(err) != nil {
		return err
	}

	r.K8S.EventRecorder.Eventf(crRef, corev1.EventTypeNormal, k8s.CreateSucceeded,
		"Created FIO Job: %v @ Namespace: %v", job.Name, job.Namespace)
	return nil
}

func (r *FioReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&perfv1alpha1.Fio{}).
		Complete(r)
}
