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
	corev1 "k8s.io/api/core/v1"

	ksapi "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("Client Pod", func() {
	Describe("created from CR", func() {
		var cr ksapi.Iperf3
		var pod *corev1.Pod

		BeforeEach(func() {
			tolerationSeconds := int64(17)
			cr = ksapi.Iperf3{
				Spec: ksapi.Iperf3Spec{
					Image: ksapi.ImageSpec{
						Name:       "foo",
						PullPolicy: "Always",
						PullSecret: "pull-secret",
					},

					ClientConfiguration: ksapi.Iperf3ConfigurationSpec{
						CmdLineArgs: "--testing --things",
						HostNetwork: true,
						PodConfigurationSpec: ksapi.PodConfigurationSpec{
							PodLabels: map[string]string{"labels": "are", "really": "useful"},
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
						},
					},
				},
			}
			pod = NewClientPod(&cr)
		})

		Context("with Image details specified", func() {
			It("should match on Image.Name", func() {
				Expect(pod.Spec.Containers[0].Image).To(
					Equal(cr.Spec.Image.Name))
			})
			It("should match on Image.PullPolicy", func() {
				Expect(pod.Spec.Containers[0].ImagePullPolicy).To(
					Equal(corev1.PullPolicy(cr.Spec.Image.PullPolicy)))
			})
			It("should match on Image.PullSecret", func() {
				Expect(pod.Spec.ImagePullSecrets[0].Name).To(
					Equal(cr.Spec.Image.PullSecret))
			})
		})

		Context("with default settings", func() {
			It("--client mode is enabled", func() {
				Expect(pod.Spec.Containers[0].Args).To(
					ContainElement("--client"))
			})
			It("port value is correctly specified", func() {
				service := NewServerService(&cr)
				servicePort := strconv.Itoa(int(service.Spec.Ports[0].Port))
				Expect(strings.Join(pod.Spec.Containers[0].Args, " ")).To(
					ContainSubstring("--port " + servicePort))
			})
			It("should not contain --udp flag", func() {
				Expect(pod.Spec.Containers[0].Args).NotTo(
					ContainElement("--udp"))
			})
		})

		Context("with UDP mode specified", func() {
			cr.Spec.UDP = true
			pod := NewClientPod(&cr)
			It("should contain --udp flag in iperf args", func() {
				Expect(pod.Spec.Containers[0].Args).To(
					ContainElement("--udp"))
			})
		})

		Context("with podLabels specified", func() {
			It("should contain all podLabels", func() {
				for k, v := range cr.Spec.ClientConfiguration.PodLabels {
					Expect(pod.ObjectMeta.Labels).To(
						HaveKeyWithValue(k, v))
				}
			})
		})

		Context("with podAffinity specified", func() {
			It("should match with Affinity", func() {
				Expect(pod.Spec.Affinity).To(
					Equal(&cr.Spec.ClientConfiguration.PodScheduling.Affinity))
			})
			It("should match with Tolerations", func() {
				Expect(pod.Spec.Tolerations).To(
					Equal(cr.Spec.ClientConfiguration.PodScheduling.Tolerations))
			})
			It("should match with NodeSelector", func() {
				Expect(pod.Spec.NodeSelector).To(
					Equal(cr.Spec.ClientConfiguration.PodScheduling.NodeSelector))
			})
			It("should match with NodeName", func() {
				Expect(pod.Spec.NodeName).To(
					Equal(cr.Spec.ClientConfiguration.PodScheduling.NodeName))
			})
		})

		Context("with HostNetwork specified", func() {
			It("should match with HostNetwork", func() {
				Expect(pod.Spec.HostNetwork).To(
					Equal(cr.Spec.ClientConfiguration.HostNetwork))
			})
		})

		Context("with connectivity to service", func() {
			It("should not match service name", func() {
				service := NewServerService(&cr)
				Expect(pod.ObjectMeta.Name).NotTo(
					Equal(service.ObjectMeta.Name))
			})
			It("should target the server service", func() {
				service := NewServerService(&cr)
				Expect(strings.Join(pod.Spec.Containers[0].Args, " ")).To(
					ContainSubstring("--client " + service.ObjectMeta.Name))
			})
		})
	})
})
