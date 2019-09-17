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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("perf job", func() {

	var objectMeta metav1.ObjectMeta
	var imageSpec perfv1alpha1.ImageSpec
	var podConfig perfv1alpha1.PodConfigurationSpec
	var job *batchv1.Job

	BeforeEach(func() {
		objectMeta = metav1.ObjectMeta{
			Name:      "test-object",
			Namespace: "test-namespace",
		}
		imageSpec = perfv1alpha1.ImageSpec{
			Name:       "test-image",
			PullPolicy: "Always",
			PullSecret: "my-secret",
		}
		tolerationSeconds := int64(17)
		podConfig = perfv1alpha1.PodConfigurationSpec{
			PodLabels: map[string]string{"this is": "an awesome test"},
			PodScheduling: perfv1alpha1.PodSchedulingSpec{
				Affinity: &corev1.Affinity{
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
				Tolerations: []corev1.Toleration{
					{
						Key:               "genetic-code",
						Operator:          corev1.TolerationOperator(corev1.TolerationOpExists),
						Value:             "distressed",
						Effect:            corev1.TaintEffect(corev1.TaintEffectNoExecute),
						TolerationSeconds: &tolerationSeconds,
					},
				},
				NodeSelector: map[string]string{
					"atomized": "spiral",
				},
				NodeName: "energy-spike-07",
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("500m"),
					corev1.ResourceMemory: resource.MustParse("5Gi"),
				},
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("1G"),
					corev1.ResourceMemory: resource.MustParse("10Gi"),
				},
			},
		}

		job = NewPerfJob(objectMeta, "test-app", imageSpec, podConfig)
	})

	Describe("NewPerfJob", func() {

		Context("ObjectMeta", func() {
			It("should match on Name", func() {
				Expect(job.ObjectMeta.Name).To(
					Equal(objectMeta.Name))
			})
			It("should match on Namespace", func() {
				Expect(job.ObjectMeta.Namespace).To(
					Equal(objectMeta.Namespace))
			})
		})

		Context("Image", func() {
			It("should match on Name", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Image).To(
					Equal(imageSpec.Name))
			})
			It("should match on PullPolicy", func() {
				Expect(job.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(
					Equal(corev1.PullPolicy(imageSpec.PullPolicy)))
			})
			It("should match on PullSecret", func() {
				Expect(job.Spec.Template.Spec.ImagePullSecrets[0].Name).To(
					Equal(imageSpec.PullSecret))
			})
		})

		It("should contain all podLabels", func() {
			for k, v := range podConfig.PodLabels {
				Expect(job.Spec.Template.ObjectMeta.Labels).To(
					HaveKeyWithValue(k, v))
			}
		})

		It("should match with Affinity", func() {
			Expect(job.Spec.Template.Spec.Affinity).To(
				Equal(podConfig.PodScheduling.Affinity))
		})
		It("should match with Tolerations", func() {
			Expect(job.Spec.Template.Spec.Tolerations).To(
				Equal(podConfig.PodScheduling.Tolerations))
		})
		It("should match with NodeSelector", func() {
			Expect(job.Spec.Template.Spec.NodeSelector).To(
				Equal(podConfig.PodScheduling.NodeSelector))
		})
		It("should match with NodeName", func() {
			Expect(job.Spec.Template.Spec.NodeName).To(
				Equal(podConfig.PodScheduling.NodeName))
		})

		Context("Resources", func() {
			It("should request the given CPU", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu()).To(
					BeEquivalentTo(podConfig.Resources.Requests.Cpu()))
			})
			It("should request the given memory", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Resources.Requests.Memory()).To(
					BeEquivalentTo(podConfig.Resources.Requests.Memory()))
			})
			It("should limit to the given CPU", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu()).To(
					BeEquivalentTo(podConfig.Resources.Limits.Cpu()))
			})
			It("should limit to the given memory", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Resources.Limits.Memory()).To(
					BeEquivalentTo(podConfig.Resources.Limits.Memory()))
			})
		})
	})
})
