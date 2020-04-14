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

package kafkabench

import (
	"context"
	"github.com/go-logr/logr"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

// KafkaBenchReconciler reconciles a KafkaBench object
type KafkaBenchReconciler struct {
	K8S k8s.Access
	Log logr.Logger
}

// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=kafkabenches,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=perf.kubestone.xridge.io,resources=kafkabenches/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;create;delete

func (r *KafkaBenchReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("kafkabench", req.NamespacedName)

	var cr perfv1alpha1.KafkaBench

	if err := r.K8S.Client.Get(ctx, req.NamespacedName, &cr); err != nil {
		return ctrl.Result{}, k8s.IgnoreNotFound(err)
	}

	// If its already completed then return
	if cr.Status.Completed {
		return ctrl.Result{}, nil
	}

	// Set status to running
	cr.Status.Running = true
	if err := r.K8S.Client.Status().Update(ctx, &cr); err != nil {
		return ctrl.Result{}, err
	}

	// Create new jobs for each test
	var jobs []*batchv1.Job
	for _, testSpec := range cr.Spec.Tests {
		result, err, kjobs := r.ProcessKafkaTest(cr, testSpec, ctx)
		if err != nil {
			return result, err
		}

		jobs = append(jobs, kjobs...)
	}

	// Check all the job statuses
	for _, job := range jobs {
		jobFinished, err := r.K8S.IsJobFinished(types.NamespacedName{
			Namespace: cr.Namespace,
			Name:      job.Name,
		})

		if err != nil {
			return ctrl.Result{}, err
		}

		if !jobFinished {
			// Wait for the job to be completed
			return ctrl.Result{Requeue: true}, nil
		}

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

func (r *KafkaBenchReconciler) ProcessKafkaTest(cr perfv1alpha1.KafkaBench, testSpec perfv1alpha1.KafkaTestSpec, ctx context.Context) (ctrl.Result, error, []*batchv1.Job) {

	// Create producer job
	producerJob := NewProducerJob(&cr, &testSpec)
	if err := r.K8S.CreateWithReference(ctx, producerJob, &cr); err != nil {
		return ctrl.Result{}, err, nil
	}

	// Create consumer job
	consumerJob := NewConsumerJob(&cr, &testSpec)
	if err := r.K8S.CreateWithReference(ctx, consumerJob, &cr); err != nil {
		return ctrl.Result{}, err, nil
	}

	return ctrl.Result{}, nil, []*batchv1.Job{consumerJob, producerJob}
}

func (r *KafkaBenchReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&perfv1alpha1.KafkaBench{}).
		Complete(r)
}

func AddPodAffinity(job *batchv1.Job, jobName string) {
	affinity := corev1.WeightedPodAffinityTerm{
		Weight: 1,
		PodAffinityTerm: corev1.PodAffinityTerm{
			LabelSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "kubestone.xridge.io/app",
						Operator: "In",
						Values:   []string{"kafkabench"},
					},
					{
						Key:      "kubestone.xridge.io/cr-name",
						Operator: "In",
						Values:   []string{jobName},
					},
				},
			},
			TopologyKey: "kubernetes.io/hostname",
		},
	}

	if job.Spec.Template.Spec.Affinity.PodAffinity == nil {
		job.Spec.Template.Spec.Affinity.PodAffinity = &corev1.PodAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution:  nil,
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{affinity},
		}
	} else {
		job.Spec.Template.Spec.Affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(
			job.Spec.Template.Spec.Affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution,
			affinity,
		)
	}
}
