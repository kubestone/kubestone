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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("pgbench job", func() {
	Describe("NewJob", func() {
		var memResource, _ = resource.ParseQuantity("4Gi")

		cr := perfv1alpha1.Pgbench{
			Spec: perfv1alpha1.PgbenchSpec{
				Image: perfv1alpha1.ImageSpec{
					Name: "postgres:latest",
				},
				Postgres: perfv1alpha1.PostgresSpec{
					Host:     "postgres",
					Port:     5432,
					User:     "admin",
					Password: "admin",
					Database: "benchdb",
				},
				PodConfig: perfv1alpha1.PodConfigurationSpec{
					Annotations:   nil,
					PodLabels:     nil,
					PodScheduling: perfv1alpha1.PodSchedulingSpec{},
					Resources: corev1.ResourceRequirements{
						Limits: nil,
						Requests: map[corev1.ResourceName]resource.Quantity{
							corev1.ResourceName("memory"): memResource,
						},
					},
				},
				InitArgs: "-s 5",
				Args:     "-t 100",
			},
		}
		job := NewJob(&cr)

		Context("both of the containers should have all the postgres env vars", func() {
			containers := job.Spec.Template.Spec.Containers
			containers = append(containers, job.Spec.Template.Spec.InitContainers...)
			for _, container := range containers {
				It("should have a host env var", func() {
					Expect(container.Env).To(
						ContainElement(corev1.EnvVar{
							Name:  "PGHOST",
							Value: cr.Spec.Postgres.Host}))
				})
				It("should have a port env var", func() {
					Expect(container.Env).To(
						ContainElement(corev1.EnvVar{
							Name:  "PGPORT",
							Value: "5432"}))
				})
				It("should have a user env var", func() {
					Expect(container.Env).To(
						ContainElement(corev1.EnvVar{
							Name:  "PGUSER",
							Value: cr.Spec.Postgres.User}))
				})
				It("should have a password env var", func() {
					Expect(container.Env).To(
						ContainElement(corev1.EnvVar{
							Name:  "PGPASSWORD",
							Value: cr.Spec.Postgres.Password}))
				})
				It("should have a database env var", func() {
					Expect(container.Env).To(
						ContainElement(corev1.EnvVar{
							Name:  "PGDATABASE",
							Value: cr.Spec.Postgres.Database}))
				})
			}
		})

		It("should run 'pgbench -i' in the init container", func() {
			Expect(job.Spec.Template.Spec.InitContainers[0].Command).To(
				Equal([]string{"pgbench", "-i"}),
			)
		})

		It("should have the given init args", func() {
			Expect(job.Spec.Template.Spec.InitContainers[0].Args).To(
				ContainElement("-s"))
			Expect(job.Spec.Template.Spec.InitContainers[0].Args).To(
				ContainElement("5"))
		})

		It("should have the given args", func() {
			Expect(job.Spec.Template.Spec.Containers[0].Args).To(
				ContainElement("-t"))
			Expect(job.Spec.Template.Spec.Containers[0].Args).To(
				ContainElement("100"))
		})

		Context("init container", func() {
			It("should have the requested resource", func() {
				Expect(job.Spec.Template.Spec.InitContainers[0].Resources).ToNot(Equal(corev1.ResourceRequirements{}))
			})
		})

	})
})
