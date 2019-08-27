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
package fio

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ksapi "github.com/xridge/kubestone/api/v1alpha1"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("fio job", func() {
	Describe("cr with cmd args", func() {
		var cr perfv1alpha1.Fio
		var job *batchv1.Job

		BeforeEach(func() {
			cr = perfv1alpha1.Fio{
				Spec: perfv1alpha1.FioSpec{
					Image: perfv1alpha1.ImageSpec{
						Name:       "xridge/fio:test",
						PullPolicy: "Always",
						PullSecret: "a-pull-secret",
					},
					CmdLineArgs: "--name=randwrite --iodepth=1 --rw=randwrite --bs=4m --direct=1 --size=256M --numjobs=1",
					PodConfig: ksapi.PodConfigurationSpec{
						PodLabels: map[string]string{"labels": "are", "still": "useful"},
						PodScheduling: ksapi.PodSchedulingSpec{
							Affinity: corev1.Affinity{
								NodeAffinity: &corev1.NodeAffinity{
									RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
										NodeSelectorTerms: []corev1.NodeSelectorTerm{
											{
												MatchExpressions: []corev1.NodeSelectorRequirement{
													{
														Key:      "mutated",
														Operator: corev1.NodeSelectorOperator(corev1.NodeSelectorOpIn),
														Values:   []string{"nano-virus"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}
			configMap := corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{Name: "cm"},
			}
			job = NewJob(&cr, &configMap, nil)
		})

		Context("with Image details specified", func() {
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

		Context("with command line args specified", func() {
			It("should have the same args", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--name=randwrite"))
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--iodepth=1"))
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--size=256M"))
			})
		})

		Context("with podAffinity specified", func() {
			It("should match with Affinity", func() {
				Expect(job.Spec.Template.Spec.Affinity).To(
					Equal(&cr.Spec.PodConfig.PodScheduling.Affinity))
			})
			It("should match with Tolerations", func() {
				Expect(job.Spec.Template.Spec.Tolerations).To(
					Equal(cr.Spec.PodConfig.PodScheduling.Tolerations))
			})
			It("should match with NodeSelector", func() {
				Expect(job.Spec.Template.Spec.NodeSelector).To(
					Equal(cr.Spec.PodConfig.PodScheduling.NodeSelector))
			})
			It("should match with NodeName", func() {
				Expect(job.Spec.Template.Spec.NodeName).To(
					Equal(cr.Spec.PodConfig.PodScheduling.NodeName))
			})
		})
	})

	Describe("cr with builtin job files and volume", func() {
		var cr perfv1alpha1.Fio
		var pvcName string
		var job *batchv1.Job

		BeforeEach(func() {
			pvcName = "test-pvc"
			cr = perfv1alpha1.Fio{
				Spec: perfv1alpha1.FioSpec{
					Image: perfv1alpha1.ImageSpec{
						Name:       "xridge/fio:test",
						PullPolicy: "IfNotPresent",
					},
					BuiltinJobFiles: []string{"/jobs/rand-read.fio"},
					Volume: &perfv1alpha1.FioVolumeSpec{
						PersistentVolumeClaimName: &pvcName,
					},
				},
			}
			configMap := corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{Name: "cm"},
			}
			job = NewJob(&cr, &configMap, &pvcName)
		})

		Context("with Image details specified", func() {
			It("should match on Image.Name", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Image).To(
					Equal(cr.Spec.Image.Name))
			})
			It("should match on Image.PullPolicy", func() {
				Expect(job.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(
					Equal(corev1.PullPolicy(cr.Spec.Image.PullPolicy)))
			})
		})

		Context("with builtin job files specified", func() {
			It("should have those files", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("/jobs/rand-read.fio"))
			})
		})

		Context("with pvc name specified", func() {
			It("we should have the pvc attached to the pod", func() {
				Expect(job.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim.ClaimName).To(
					Equal(pvcName))
			})
		})

		Context("with podLabels specified", func() {
			It("should contain all podLabels", func() {
				for key, value := range cr.Spec.PodConfig.PodLabels {
					Expect(job.Spec.Template.ObjectMeta.Labels).To(
						HaveKeyWithValue(key, value))
				}
			})
		})
	})

	Describe("cr with volume without pvc", func() {
		var cr perfv1alpha1.Fio
		var job *batchv1.Job

		BeforeEach(func() {
			cr = perfv1alpha1.Fio{
				Spec: perfv1alpha1.FioSpec{
					Image: perfv1alpha1.ImageSpec{
						Name:       "xridge/fio:test",
						PullPolicy: "IfNotPresent",
					},
					Volume: &perfv1alpha1.FioVolumeSpec{},
				},
			}
			job = NewJob(&cr, nil, nil)
		})

		Context("with Volume withohout pvc", func() {
			It("should create an emptydir", func() {
				Expect(job.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim).To(
					BeNil())
			})
		})
	})
})
