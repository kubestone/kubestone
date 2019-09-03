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
	corev1 "k8s.io/api/core/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("drill configmap", func() {
	Describe("cr with files for benchmarkVolumes", func() {
		var cr perfv1alpha1.Drill
		var configMap *corev1.ConfigMap

		BeforeEach(func() {
			cr = perfv1alpha1.Drill{
				Spec: perfv1alpha1.DrillSpec{
					BenchmarksVolume: map[string]string{
						"file-1.yml": "content-1",
						"file-2.yml": "content-2",
					},
				},
			}
			configMap = NewConfigMap(&cr)
		})

		Context("with benchmark files specified", func() {
			It("should have them in the configmap", func() {
				Expect(configMap.Data["file-1.yml"]).To(
					Equal(cr.Spec.BenchmarksVolume["file-1.yml"]))
				Expect(configMap.Data["file-2.yml"]).To(
					Equal(cr.Spec.BenchmarksVolume["file-2.yml"]))
			})
		})
	})
})
