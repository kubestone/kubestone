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
package iperf3

import (
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	ksapi "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("Server Deployment", func() {
	Describe("created from CR", func() {
		var cr ksapi.Iperf3
		var deployment *appsv1.Deployment

		BeforeEach(func() {
			tolerationSeconds := int64(17)
			cr = ksapi.Iperf3{
				Spec: ksapi.Iperf3Spec{
					Image: ksapi.ImageSpec{
						Name:       "foo",
						PullPolicy: "Always",
						PullSecret: "pull-secret",
					},

					ServerConfiguration: ksapi.Iperf3ConfigurationSpec{
						CmdLineArgs: "--testing --things",
						HostNetwork: true,
						PodConfigurationSpec: ksapi.PodConfigurationSpec{
							PodLabels: map[string]string{"labels": "are", "really": "useful"},
							PodScheduling: ksapi.PodSchedulingSpec{
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
						},
					},
				},
			}
			deployment = NewServerDeployment(&cr)
		})

		Context("with Image details specified", func() {
			It("should match on Image.Name", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Image).To(
					Equal(cr.Spec.Image.Name))
			})
			It("should match on Image.PullPolicy", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(
					Equal(corev1.PullPolicy(cr.Spec.Image.PullPolicy)))
			})
			It("should match on Image.PullSecret", func() {
				Expect(deployment.Spec.Template.Spec.ImagePullSecrets[0].Name).To(
					Equal(cr.Spec.Image.PullSecret))
			})
		})

		Context("with default settings", func() {
			It("--server mode is enabled", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--server"))
			})
			It("--port's value is specified", func() {
				Expect(strings.Join(deployment.Spec.Template.Spec.Containers[0].Args, " ")).To(
					ContainSubstring("--port " + strconv.Itoa(Iperf3ServerPort)))
			})
			It("should not contain --udp flag", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Args).NotTo(
					ContainElement("--udp"))
			})
		})

		Context("with UDP mode specified", func() {
			cr.Spec.UDP = true
			deployment := NewServerDeployment(&cr)
			It("should contain --udp flag in iperf args", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--udp"))
			})
		})

		Context("with cmdLineArgs specified", func() {
			It("--testing mode is set", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--testing"))
			})
		})

		Context("with podLabels specified", func() {
			It("should contain all podLabels", func() {
				for k, v := range cr.Spec.ServerConfiguration.PodLabels {
					Expect(deployment.Spec.Template.ObjectMeta.Labels).To(
						HaveKeyWithValue(k, v))
				}
			})
		})

		Context("with podAffinity specified", func() {
			It("should match with Affinity", func() {
				Expect(deployment.Spec.Template.Spec.Affinity).To(
					Equal(cr.Spec.ServerConfiguration.PodScheduling.Affinity))
			})
			It("should match with Tolerations", func() {
				Expect(deployment.Spec.Template.Spec.Tolerations).To(
					Equal(cr.Spec.ServerConfiguration.PodScheduling.Tolerations))
			})
			It("should match with NodeSelector", func() {
				Expect(deployment.Spec.Template.Spec.NodeSelector).To(
					Equal(cr.Spec.ServerConfiguration.PodScheduling.NodeSelector))
			})
			It("should match with NodeName", func() {
				Expect(deployment.Spec.Template.Spec.NodeName).To(
					Equal(cr.Spec.ServerConfiguration.PodScheduling.NodeName))
			})
		})

		Context("with HostNetwork specified", func() {
			It("should match with HostNetwork", func() {
				Expect(deployment.Spec.Template.Spec.HostNetwork).To(
					Equal(cr.Spec.ServerConfiguration.HostNetwork))
			})
		})

		Context("with resources specified", func() {
			It("should request the given CPU", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu()).To(
					BeEquivalentTo(cr.Spec.ServerConfiguration.Resources.Requests.Cpu()))
			})
			It("should request the given memory", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory()).To(
					BeEquivalentTo(cr.Spec.ServerConfiguration.Resources.Requests.Memory()))
			})
			It("should limit to the given CPU", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu()).To(
					BeEquivalentTo(cr.Spec.ServerConfiguration.Resources.Limits.Cpu()))
			})
			It("should limit to the given memory", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory()).To(
					BeEquivalentTo(cr.Spec.ServerConfiguration.Resources.Limits.Memory()))
			})
		})
	})
})
