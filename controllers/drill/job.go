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
	"errors"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
)

const (
	benchmarksDir = "/benchmarks"
	drill         = "drill"
)

// NewJob creates a fio benchmark job
func NewJob(cr *perfv1alpha1.Drill, configMap *corev1.ConfigMap) *batchv1.Job {
	objectMeta := metav1.ObjectMeta{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}

	cmdLineArgs := fmt.Sprintf("%s --benchmark %s", cr.Spec.Options, cr.Spec.BenchmarkFile)
	command := fmt.Sprintf("cd %s && %s %s", benchmarksDir, drill, cmdLineArgs)

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
			Name:      "benchmarks",
			MountPath: benchmarksDir,
		},
	}

	job := k8s.NewPerfJob(objectMeta, "fio", cr.Spec.Image, cr.Spec.PodConfig)
	job.Spec.Template.Spec.Volumes = volumes
	job.Spec.Template.Spec.Containers[0].Command = []string{"/bin/sh", "-xc"}
	job.Spec.Template.Spec.Containers[0].Args = []string{command}
	job.Spec.Template.Spec.Containers[0].VolumeMounts = volumeMounts
	return job
}

// IsCrValid validates the given CR and raises error if semantic errors detected
// For drill it checks that the BenchmarkFile exists in the BenchmarksVolume map
func IsCrValid(cr *perfv1alpha1.Drill) (valid bool, err error) {
	if _, ok := cr.Spec.BenchmarksVolume[cr.Spec.BenchmarkFile]; !ok {
		return false, errors.New("BenchmarkFile does not exists in BenchmarksVolume")
	}

	return true, nil
}
