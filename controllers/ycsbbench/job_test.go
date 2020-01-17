package ycsbbench

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("ycsbbench job", func() {
	Describe("NewJob", func() {
		cr := perfv1alpha1.YcsbBench{
			Spec: perfv1alpha1.YcsbBenchSpec{
				Image: perfv1alpha1.ImageSpec{
					Name: "diamantisolutions/ycsb:latest",
				},
				Action:     "load",
				DbType:     "redis",
				Workletter: "a",
				DbArgs:     "-p redis.host=10.0.0.1",
			},
		}

		job := NewJob(&cr)

		Context("container should have all the envs", func() {
			containers := job.Spec.Template.Spec.Containers
			for _, container := range containers {
				It("should have action env var", func() {
					Expect(container.Env).To(
						ContainElement(corev1.EnvVar{
							Name:  "ACTION",
							Value: "load"}))
				})

				It("should have db_type env var", func() {
					Expect(container.Env).To(
						ContainElement(corev1.EnvVar{
							Name:  "DBTYPE",
							Value: "redis"}))
				})

				It("should have workletter env var", func() {
					Expect(container.Env).To(
						ContainElement(corev1.EnvVar{
							Name:  "WORKLETTER",
							Value: "a"}))
				})

				It("should have db_args env var", func() {
					Expect(container.Env).To(
						ContainElement(corev1.EnvVar{
							Name:  "DBARGS",
							Value: "-p redis.host=10.0.0.1"}))
				})
			}
		})
	})
})
