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

package sysbench

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/firepear/qsplit"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// NewJob creates a sysbench benchmark job
func NewJob(cr *perfv1alpha1.Sysbench) *batchv1.Job {
	labels := map[string]string{
		"app":               "sysbench",
		"kubestone-cr-name": cr.Name,
	}
	for key, value := range cr.Spec.PodLabels {
		labels[key] = value
	}

	sysbenchCmdLineArgs := []string{}
	sysbenchCmdLineArgs = append(sysbenchCmdLineArgs, qsplit.ToStrings([]byte(cr.Spec.Options))...)
	sysbenchCmdLineArgs = append(sysbenchCmdLineArgs, cr.Spec.TestName, cr.Spec.Command)

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
							Name:            "sysbench",
							Image:           cr.Spec.Image.Name,
							ImagePullPolicy: corev1.PullPolicy(cr.Spec.Image.PullPolicy),
							Args:            sysbenchCmdLineArgs,
						},
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: cr.Spec.Image.PullSecret,
						},
					},
					// TODO: add more options here eg. resource requests/limits
					RestartPolicy: corev1.RestartPolicyNever,
					Affinity:      &cr.Spec.PodScheduling.Affinity,
					Tolerations:   cr.Spec.PodScheduling.Tolerations,
					NodeSelector:  cr.Spec.PodScheduling.NodeSelector,
					NodeName:      cr.Spec.PodScheduling.NodeName,
				},
			},
			BackoffLimit: &backoffLimit,
		},
	}

	return &job
}
