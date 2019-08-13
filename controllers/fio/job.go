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

package fio

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/firepear/qsplit"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// NewJob creates a fio benchmark job
func NewJob(cr *perfv1alpha1.Fio) *batchv1.Job {
	labels := map[string]string{
		"app":               "fio",
		"kubestone-cr-name": cr.Name,
	}

	fioCmdLineArgs := []string{}

	fioCmdLineArgs = append(fioCmdLineArgs,
		qsplit.ToStrings([]byte(cr.Spec.CmdLineArgs))...)

	for _, builtinJobFile := range cr.Spec.BuiltinJobFiles {
		fioCmdLineArgs = append(fioCmdLineArgs, builtinJobFile)
	}

	backoffLimit := int32(0)

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
							ImagePullPolicy: corev1.PullPolicy(cr.Spec.Image.PullPolicy),
							Args:            fioCmdLineArgs,
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backoffLimit,
		},
	}

	return &job
}

func (r *Reconciler) isJobFinished(cr *perfv1alpha1.Fio) (finished bool, err error) {
	// TODO: Move this to k8s.client
	job, err := r.K8S.Clientset.BatchV1().Jobs(cr.Namespace).Get(cr.Name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	finished = job.Status.Succeeded+job.Status.Failed > 0
	return finished, nil
}
