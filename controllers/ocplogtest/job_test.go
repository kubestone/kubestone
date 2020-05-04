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
package ocplogtest

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("ocplogtest job", func() {
	Describe("NewJob", func() {
		cr := perfv1alpha1.OcpLogtest{
			Spec: perfv1alpha1.OcpLogtestSpec{
				Image: perfv1alpha1.ImageSpec{
					Name: "quay.io/mffiedler/ocp-logtest:latest",
				},
				LineLength: 1024,
				NumLines:   300000,
				Rate:       60000,
				FixedLine:  true,
			},
		}
		job := NewJob(&cr)

		It("should run 'python' in the job container", func() {
			Expect(job.Spec.Template.Spec.Containers[0].Command).To(
				Equal([]string{"python"}),
			)
		})

		It("should have the translated args", func() {
			Expect(job.Spec.Template.Spec.Containers[0].Args).To(
				ContainElement("ocp_logtest.py"))

			Expect(job.Spec.Template.Spec.Containers[0].Args).To(
				ContainElement("--line-length"))
			Expect(job.Spec.Template.Spec.Containers[0].Args).To(
				ContainElement("1024"))

			Expect(job.Spec.Template.Spec.Containers[0].Args).To(
				ContainElement("--num-lines"))
			Expect(job.Spec.Template.Spec.Containers[0].Args).To(
				ContainElement("300000"))

			Expect(job.Spec.Template.Spec.Containers[0].Args).To(
				ContainElement("--rate"))
			Expect(job.Spec.Template.Spec.Containers[0].Args).To(
				ContainElement("60000"))

			Expect(job.Spec.Template.Spec.Containers[0].Args).To(
				ContainElement("--fixed-line"))
		})
	})
})
