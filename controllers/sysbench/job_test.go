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
package sysbench

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("sysbench job", func() {
	Describe("cr with cmd args", func() {
		var cr perfv1alpha1.Sysbench
		var job *batchv1.Job

		BeforeEach(func() {
			cr = perfv1alpha1.Sysbench{
				Spec: perfv1alpha1.SysbenchSpec{
					Image: perfv1alpha1.ImageSpec{
						Name: "xridge/sysbench:test",
					},
					Options:  "--threads=2 --time=20",
					TestName: "cpu",
					Command:  "run",
				},
			}
			job = NewJob(&cr)
		})

		Context("with command line arguments", func() {
			It("should have the same options", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--threads=2"))
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--time=20"))
			})
			It("should have the same testName", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("cpu"))
			})
			It("should have the same command", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("run"))
			})
		})
	})
})
