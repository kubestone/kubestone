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
package s3bench

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("s3bench job", func() {
	Describe("cr with cmd args", func() {
		var cr perfv1alpha1.S3Bench
		var job *batchv1.Job

		BeforeEach(func() {
			cr = perfv1alpha1.S3Bench{Spec: perfv1alpha1.S3BenchSpec{
				Mode: "mixed",
				Host: "minio-test.minio.svc.sol1.diamanti.com:9000",
				S3BenchOptions: perfv1alpha1.S3BenchOptions{
					AccessKey: "AKIAIOSFODNN7EXAMPLE",
					SecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				},
				S3ObjectOptions:          perfv1alpha1.S3ObjectOptions{},
				S3AutoTermOptions:        perfv1alpha1.S3AutoTermOptions{},
				S3AnalysisOptions:        perfv1alpha1.S3AnalysisOptions{},
				MixedDistributionOptions: perfv1alpha1.MixedDistributionOptions{},
			}}
			job = NewJob(&cr)
		})

		Context("with command line arguments", func() {
			It("should have the same options", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("mixed"))
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("minio-test.minio.svc.sol1.diamanti.com:9000"))
			})
		})
	})
})
