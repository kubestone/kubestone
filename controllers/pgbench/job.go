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

package pgbench

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/firepear/qsplit"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
)

// NewJob creates a new pgbench job
func NewJob(cr *perfv1alpha1.Pgbench) *batchv1.Job {
	objectMeta := metav1.ObjectMeta{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}

	env := []corev1.EnvVar{
		{Name: "PGHOST", Value: cr.Spec.Postgres.Host},
		{Name: "PGPORT", Value: fmt.Sprintf("%d", cr.Spec.Postgres.Port)},
		{Name: "PGUSER", Value: cr.Spec.Postgres.User},
		{Name: "PGPASSWORD", Value: cr.Spec.Postgres.Password},
		{Name: "PGDATABASE", Value: cr.Spec.Postgres.Database},
	}

	initContainer := corev1.Container{
		Name:            "pgbench-init",
		Image:           cr.Spec.Image.Name,
		ImagePullPolicy: corev1.PullPolicy(cr.Spec.Image.PullPolicy),
		Command:         []string{"pgbench", "-i"},
		Args:            qsplit.ToStrings([]byte(cr.Spec.InitArgs)),
		Env:             env,
		Resources:       cr.Spec.PodConfig.Resources,
	}

	job := k8s.NewPerfJob(objectMeta, "pgbench", cr.Spec.Image, cr.Spec.PodConfig)
	job.Spec.Template.Spec.InitContainers = append(
		job.Spec.Template.Spec.InitContainers, initContainer)
	job.Spec.Template.Spec.Containers[0].Args = qsplit.ToStrings([]byte(cr.Spec.Args))
	job.Spec.Template.Spec.Containers[0].Env = env
	return job
}
