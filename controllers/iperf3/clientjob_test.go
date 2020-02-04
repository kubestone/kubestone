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
	batchv1 "k8s.io/api/batch/v1"

	ksapi "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("Client Pod", func() {
	Describe("created from CR", func() {
		var cr ksapi.Iperf3
		var job *batchv1.Job

		BeforeEach(func() {
			cr = ksapi.Iperf3{
				Spec: ksapi.Iperf3Spec{
					Image: ksapi.ImageSpec{
						Name: "foo",
					},
					ClientConfiguration: ksapi.Iperf3ConfigurationSpec{
						CmdLineArgs: "--testing --things",
						HostNetwork: true,
						PodConfigurationSpec: ksapi.PodConfigurationSpec{
							Annotations: map[string]string{
								"annotation_one": "value_one",
							},
						},
					},
				},
			}
			job = NewClientJob(&cr)
		})

		Context("with default settings", func() {
			It("--client mode is enabled", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--client"))
			})
			It("port value is correctly specified", func() {
				service := NewServerService(&cr)
				servicePort := strconv.Itoa(int(service.Spec.Ports[0].Port))
				Expect(strings.Join(job.Spec.Template.Spec.Containers[0].Args, " ")).To(
					ContainSubstring("--port " + servicePort))
			})
			It("should not contain --udp flag", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).NotTo(
					ContainElement("--udp"))
			})
		})

		Context("with cmdLineArgs specified", func() {
			It("--testing mode is set", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--testing"))
			})
		})

		Context("with UDP mode specified", func() {
			cr.Spec.UDP = true
			job := NewClientJob(&cr)
			It("should contain --udp flag in iperf args", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--udp"))
			})
		})

		Context("with HostNetwork specified", func() {
			It("should match with HostNetwork", func() {
				Expect(job.Spec.Template.Spec.HostNetwork).To(
					Equal(cr.Spec.ClientConfiguration.HostNetwork))
			})
		})

		Context("with connectivity to service", func() {
			It("should not match service name", func() {
				service := NewServerService(&cr)
				Expect(job.ObjectMeta.Name).NotTo(
					Equal(service.ObjectMeta.Name))
			})
			It("should target the server service", func() {
				service := NewServerService(&cr)
				Expect(strings.Join(job.Spec.Template.Spec.Containers[0].Args, " ")).To(
					ContainSubstring("--client " + service.ObjectMeta.Name))
			})
		})

		Context("by default", func() {
			defaultBackoffLimit := int32(6)
			It("should retry 6 times", func() {
				Expect(job.Spec.BackoffLimit).To(
					Equal(&defaultBackoffLimit))
			})
		})

		Context("with added annotations", func() {
			It("should contain pod annotations", func() {
				Expect(job.ObjectMeta.Annotations).To(HaveKey("annotation_one"))
			})
		})
	})
})
