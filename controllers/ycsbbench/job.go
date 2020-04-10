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

package ycsbbench

import (
	"fmt"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

func formatArgs(cr *perfv1alpha1.YcsbBench) []string {
	args := []string{
		cr.Spec.Database,
		"-P", fmt.Sprintf("workloads/workload%s", cr.Spec.Workload),
	}

	if cr.Spec.Options != (perfv1alpha1.YcsbBenchOptions{}) {
		if cr.Spec.Options.Threadcount > 0 {
			args = append(args, "-threads", strconv.Itoa(cr.Spec.Options.Threadcount))
		}

		if cr.Spec.Options.Target > 0 {
			args = append(args, "-target", strconv.Itoa(cr.Spec.Options.Target))
		}
	}

	for key, val := range cr.Spec.Properties {
		args = append(args, "-p", fmt.Sprintf("%s=%s", key, val))
	}
	return args
}

func NewJob(cr *perfv1alpha1.YcsbBench) *batchv1.Job {
	objectMeta := metav1.ObjectMeta{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}

	args := formatArgs(cr)
	initContainer := corev1.Container{
		Name:            "ycsbbench-load",
		Image:           cr.Spec.Image.Name,
		ImagePullPolicy: corev1.PullPolicy(cr.Spec.Image.PullPolicy),
		Command:         []string{"./bin/ycsb", "load"},
		Args:            args,
	}
	// append([]string{"./bin/ycsb", "load", args},
	job := k8s.NewPerfJob(objectMeta, "ycsbbench", cr.Spec.Image, cr.Spec.PodConfig)
	job.Spec.Template.Spec.InitContainers = append(
		job.Spec.Template.Spec.InitContainers, initContainer)

	job.Spec.Template.Spec.Containers[0].Command = []string{"./bin/ycsb", "run"}
	job.Spec.Template.Spec.Containers[0].Args = args

	return job
}
