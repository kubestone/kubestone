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

package jmeter

import (
	"errors"
	"fmt"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/firepear/qsplit"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
)

const (
	testsDir      = "/jmeter-plan-tests"
	propertiesDir = "/jmeter-properties"
	reportsDir    = "/jmeter-reports"
)

// NewJob creates a new jmeter job
func NewJob(cr *perfv1alpha1.JMeter, planTestConfigMap, propertiesConfigMap *corev1.ConfigMap) *batchv1.Job {
	objectMeta := metav1.ObjectMeta{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}

	volumes := generateVolumes(cr, planTestConfigMap, propertiesConfigMap)
	volumeMounts := generateVolumeMounts(cr, planTestConfigMap, propertiesConfigMap)

	args := qsplit.ToStrings([]byte(cr.Spec.Args))
	args = append(args,
		"-t",
		fmt.Sprintf("%s/%s", testsDir, cr.Spec.TestName),
		"-o",
		reportsDir,
	)

	if propertiesConfigMap != nil {
		args = append(args, "-p", fmt.Sprintf("%s/%s", propertiesDir, cr.Spec.PropsName))
	}

	job := k8s.NewPerfJob(objectMeta, "jmeter", cr.Spec.Image, cr.Spec.Configuration)
	job.Spec.Completions = cr.Spec.JobConfig.Completions
	job.Spec.Parallelism = cr.Spec.JobConfig.Parallelism
	job.Spec.Template.Spec.Volumes = volumes
	job.Spec.Template.Spec.Containers[0].Args = args
	job.Spec.Template.Spec.Containers[0].Command = qsplit.ToStrings([]byte(cr.Spec.Command))
	job.Spec.Template.Spec.Containers[0].VolumeMounts = volumeMounts
	return job
}

func generateVolumes(cr *perfv1alpha1.JMeter, planTestConfigMap, propertiesConfigMap *corev1.ConfigMap) []corev1.Volume {
	volumes := []corev1.Volume{
		corev1.Volume{
			Name: "plans",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: planTestConfigMap.Name,
					},
				},
			},
		},
		corev1.Volume{
			Name:         "reports",
			VolumeSource: cr.Spec.Volume.VolumeSource,
		},
	}

	if propertiesConfigMap != nil {
		propertiesVolume := corev1.Volume{
			Name: "properties",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: propertiesConfigMap.Name,
					},
				},
			},
		}

		volumes = append(volumes, propertiesVolume)
	}

	return volumes
}

func generateVolumeMounts(cr *perfv1alpha1.JMeter, planTestConfigMap, propertiesConfigMap *corev1.ConfigMap) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		corev1.VolumeMount{
			Name:      "plans",
			MountPath: testsDir,
		},
		corev1.VolumeMount{
			Name:      "reports",
			MountPath: reportsDir,
		},
	}

	if propertiesConfigMap != nil {
		propertiesVolumeMount := corev1.VolumeMount{
			Name:      "properties",
			MountPath: propertiesDir,
		}
		volumeMounts = append(volumeMounts, propertiesVolumeMount)
	}

	return volumeMounts
}

// IsCrValid validates the given CR and raises error if semantic errors detected
// For jmeter it checks that the plan test is valid
func IsCrValid(cr *perfv1alpha1.JMeter) (valid bool, err error) {
	if len(cr.Spec.TestName) == 0 {
		return false, errors.New("You need to specify the TestName")
	}

	if strings.Contains(cr.Spec.Args, "-t") {
		return false, fmt.Errorf("You can't specify the flag '-t'")
	}

	if strings.Contains(cr.Spec.Args, "-o") {
		return false, fmt.Errorf("You can't specify the flag '-o'")
	}

	testName := cr.Spec.TestName
	planTest, ok := cr.Spec.PlanTest[testName]

	if !ok {
		return false, fmt.Errorf("The key '%s' is missing at spec.planTest", testName)
	}

	if planTest == "" {
		return false, fmt.Errorf("The key '%s' is empty at spec.planTest", testName)
	}

	if ok, err := cr.Spec.Volume.Validate(); !ok || err != nil {
		return false, fmt.Errorf("The volume spec is invalid: %s", err)
	}

	return true, nil
}
