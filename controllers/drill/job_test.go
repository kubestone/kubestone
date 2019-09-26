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

package drill

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("drill job", func() {
	Describe("cr with cmd args", func() {
		var cr perfv1alpha1.Drill
		var job *batchv1.Job

		BeforeEach(func() {
			cr = perfv1alpha1.Drill{
				Spec: perfv1alpha1.DrillSpec{
					Image: perfv1alpha1.ImageSpec{
						Name:       "xridge/drill:test",
						PullPolicy: "Always",
						PullSecret: "the-pull-secret",
					},
					BenchmarksVolume: map[string]string{
						"the-benchmark.yml": "benchmark content",
						"included-file.yml": "included content",
					},
					BenchmarkFile: "the-benchmark.yml",
					Options:       "--no-check-certificate --stats",
				},
			}
			configMap := corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{Name: "cm"},
			}
			job = NewJob(&cr, &configMap)
		})
		Context("with command line args specified", func() {
			It("should have the same args", func() {
				Expect(len(job.Spec.Template.Spec.Containers[0].Args)).To(Equal(1))
				args := job.Spec.Template.Spec.Containers[0].Args[0]
				Expect(args).To(ContainSubstring("--no-check-certificate"))
				Expect(args).To(ContainSubstring("--stats"))
				Expect(args).To(ContainSubstring("--benchmark"))
				Expect(args).To(ContainSubstring(cr.Spec.BenchmarkFile))
			})
		})

		Context("when existent benchmarkFile is referred", func() {
			It("CR Validation should succeed", func() {
				valid, err := IsCrValid(&cr)
				Expect(valid).To(BeTrue())
				Expect(err).To(BeNil())
			})
		})
		Context("when non-existent benchmarkFile is referred", func() {
			invalidCr := perfv1alpha1.Drill{
				Spec: perfv1alpha1.DrillSpec{
					Image: perfv1alpha1.ImageSpec{
						Name: "xridge/drill:test",
					},
					BenchmarksVolume: map[string]string{
						"the-benchmark.yml": "benchmark content",
					},
					BenchmarkFile: "non-existent-benchmark.yml",
				},
			}

			It("CR Validation should fail", func() {
				valid, err := IsCrValid(&invalidCr)
				Expect(valid).To(BeFalse())
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
