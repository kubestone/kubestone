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

package qperf

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
		var cr ksapi.Qperf
		var job *batchv1.Job

		BeforeEach(func() {
			cr = ksapi.Qperf{
				Spec: ksapi.QperfSpec{
					Image: ksapi.ImageSpec{
						Name: "foo",
					},
					Tests: []string{
						"tcp_bw",
						"tcp_lat",
					},
					ClientConfiguration: ksapi.QperfConfigurationSpec{
						HostNetwork: true,
					},
				},
			}
			job = NewClientJob(&cr)
		})

		Context("with default settings", func() {
			It("port value is correctly specified", func() {
				service := NewServerService(&cr)
				servicePort := strconv.Itoa(int(service.Spec.Ports[0].Port))
				Expect(strings.Join(job.Spec.Template.Spec.Containers[0].Args, " ")).To(
					ContainSubstring("--listen_port " + servicePort))
			})
		})

		Context("with Tests specified", func() {
			It("should appear in command line args", func() {
				length := len(job.Spec.Template.Spec.Containers[0].Args) - len(cr.Spec.Tests)
				Expect(job.Spec.Template.Spec.Containers[0].Args[length:]).To(
					Equal(cr.Spec.Tests))
			})
		})

		Context("with Options specified", func() {
			It("should appear in command line args", func() {
				cr.Spec.Options = "--option1 --option2"
				job = NewClientJob(&cr)
				Expect(strings.Join(job.Spec.Template.Spec.Containers[0].Args, " ")).To(
					ContainSubstring(cr.Spec.Options))
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
				Expect(job.Spec.Template.Spec.Containers[0].Args[0]).To(
					Equal(service.ObjectMeta.Name))
			})
		})

		Context("by default", func() {
			defaultBackoffLimit := int32(6)
			It("should retry 6 times", func() {
				Expect(job.Spec.BackoffLimit).To(
					Equal(&defaultBackoffLimit))
			})
		})
	})
})
