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

package k8s

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// NewPerfJob creates a new kubernetes job with the given arguments
func NewPerfJob(objectMeta metav1.ObjectMeta, app string, imageSpec perfv1alpha1.ImageSpec,
	podConfig perfv1alpha1.PodConfigurationSpec) *batchv1.Job {

	backoffLimit := int32(0)

	labels := map[string]string{
		"kubestone.xridge.io/app":     app,
		"kubestone.xridge.io/cr-name": objectMeta.Name,
	}
	for key, value := range podConfig.PodLabels {
		labels[key] = value
	}

	if podConfig.Annotations != nil {
		if objectMeta.Annotations != nil {
			// If its already set add the values onto the existing map
			for k, v := range podConfig.Annotations {
				objectMeta.Annotations[k] = v
			}
		} else {
			objectMeta.Annotations = podConfig.Annotations
		}
	}

	if podConfig.PodScheduling.Affinity == nil {
		podConfig.PodScheduling.Affinity = &corev1.Affinity{
			NodeAffinity: &corev1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution:  nil,
				PreferredDuringSchedulingIgnoredDuringExecution: nil,
			},
			PodAffinity: &corev1.PodAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution:  nil,
				PreferredDuringSchedulingIgnoredDuringExecution: nil,
			},
			PodAntiAffinity: &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution:  nil,
				PreferredDuringSchedulingIgnoredDuringExecution: nil,
			},
		}
	}

	job := batchv1.Job{
		ObjectMeta: objectMeta,
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: objectMeta.Annotations,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            app,
							Image:           imageSpec.Name,
							ImagePullPolicy: corev1.PullPolicy(imageSpec.PullPolicy),
							Resources:       podConfig.Resources,
						},
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						{Name: imageSpec.PullSecret},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Affinity:      podConfig.PodScheduling.Affinity,
					Tolerations:   podConfig.PodScheduling.Tolerations,
					NodeSelector:  podConfig.PodScheduling.NodeSelector,
					NodeName:      podConfig.PodScheduling.NodeName,
				},
			},
			BackoffLimit: &backoffLimit,
		},
	}

	return &job
}
