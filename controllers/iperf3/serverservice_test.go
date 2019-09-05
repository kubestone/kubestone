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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	ksapi "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("Server Service", func() {
	Describe("created from CR", func() {
		var cr ksapi.Iperf3
		var service *corev1.Service

		BeforeEach(func() {
			cr = ksapi.Iperf3{
				Spec: ksapi.Iperf3Spec{
					Image: ksapi.ImageSpec{
						Name: "foo",
					},
				},
			}
			service = NewServerService(&cr)
		})

		Context("with default settings", func() {
			It("should use TCP protocol", func() {
				Expect(service.Spec.Ports[0].Protocol).To(
					Equal(corev1.ProtocolTCP))
			})
		})

		Context("with UDP mode specified", func() {
			cr.Spec.UDP = true
			service := NewServerService(&cr)
			It("should use UDP protocol", func() {
				Expect(service.Spec.Ports[0].Protocol).To(
					Equal(corev1.ProtocolUDP))
			})
		})

		Context("crosschecked with server deployment", func() {
			service := NewServerService(&cr)
			deployment := NewServerDeployment(&cr)
			It("should match on port", func() {
				Expect(service.Spec.Ports[0].Protocol).To(
					Equal(deployment.Spec.Template.Spec.Containers[0].Ports[0].Protocol))
			})
			It("should match on selectors", func() {
				Expect(service.Spec.Selector).To(
					Equal(deployment.Spec.Template.ObjectMeta.Labels))
			})
			It("should match on namespace", func() {
				Expect(service.ObjectMeta.Namespace).To(
					Equal(deployment.ObjectMeta.Namespace))
			})
		})
	})
})
