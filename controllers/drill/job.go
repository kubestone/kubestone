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

package drill

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/firepear/qsplit"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// NewJob creates a fio benchmark job
func NewJob(cr *perfv1alpha1.Drill, configMap *corev1.ConfigMap) *batchv1.Job {
	labels := map[string]string{
		"app":               "fio",
		"kubestone-cr-name": cr.Name,
	}
	for key, value := range cr.Spec.PodConfig.PodLabels {
		labels[key] = value
	}

	cmdLineArgs := []string{}
	cmdLineArgs = append(cmdLineArgs, qsplit.ToStrings([]byte(cr.Spec.Options))...)
	cmdLineArgs = append(cmdLineArgs, "--benchmark", cr.Spec.BenchmarkFile)

	backoffLimit := int32(0)

	volumes := []corev1.Volume{
		corev1.Volume{
			Name: "benchmarks",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configMap.Name,
					},
				},
			},
		},
	}
	volumeMounts := []corev1.VolumeMount{
		corev1.VolumeMount{
			Name: configMap.Name, MountPath: "/benchmarks",
		},
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
							Name:            "drill",
							Image:           cr.Spec.Image.Name,
							ImagePullPolicy: corev1.PullPolicy(cr.Spec.Image.PullPolicy),
							Args:            cmdLineArgs,
							VolumeMounts:    volumeMounts,
						},
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: cr.Spec.Image.PullSecret,
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes:       volumes,
					Affinity:      &cr.Spec.PodConfig.PodScheduling.Affinity,
					Tolerations:   cr.Spec.PodConfig.PodScheduling.Tolerations,
					NodeSelector:  cr.Spec.PodConfig.PodScheduling.NodeSelector,
					NodeName:      cr.Spec.PodConfig.PodScheduling.NodeName,
				},
			},
			BackoffLimit: &backoffLimit,
		},
	}

	return &job
}
