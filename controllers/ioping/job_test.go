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
package ioping

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("ioping job", func() {
	Describe("with minimum parameter set", func() {
		var cr perfv1alpha1.Ioping
		var job *batchv1.Job

		BeforeEach(func() {
			cr = perfv1alpha1.Ioping{
				Spec: perfv1alpha1.IopingSpec{
					Image: perfv1alpha1.ImageSpec{
						Name:       "xridge/ioping:test",
						PullPolicy: "Always",
						PullSecret: "a-pull-secret",
					},
					Args: "-P 42",
					Volume: perfv1alpha1.VolumeSpec{
						VolumeSource: corev1.VolumeSource{
							EmptyDir: &corev1.EmptyDirVolumeSource{
								Medium: "floppy",
							},
						},
					},
				},
			}
			job = NewJob(&cr)
		})

		Context("with command line args specified", func() {
			It("should have the target dir specified", func() {
				lastElement := len(job.Spec.Template.Spec.Containers[0].Args) - 1
				Expect(job.Spec.Template.Spec.Containers[0].Args[lastElement]).To(
					Equal("/data"))
			})
			It("should have the provided args", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("-P"))
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("42"))
			})
		})

		Context("with image details specified", func() {
			It("should match on Image.Name", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Image).To(
					Equal(cr.Spec.Image.Name))
			})
			It("should match on Image.PullPolicy", func() {
				Expect(job.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(
					Equal(corev1.PullPolicy(cr.Spec.Image.PullPolicy)))
			})
			It("should match on Image.PullSecret", func() {
				Expect(job.Spec.Template.Spec.ImagePullSecrets[0].Name).To(
					Equal(cr.Spec.Image.PullSecret))
			})
		})
		Context("with volumeSource defined", func() {
			It("should match on provided volume", func() {
				Expect(job.Spec.Template.Spec.Volumes[0].VolumeSource).To(
					Equal(cr.Spec.Volume.VolumeSource))
			})
			It("should mount as /data", func() {
				Expect(job.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath).To(
					Equal("/data"))
			})
		})
	})
})
