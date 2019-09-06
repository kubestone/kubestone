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

package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xridge/kubestone/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	sysbenchSampleCR = samplesDir + "/perf_v1alpha1_sysbench.yaml"
)

var _ = Describe("sysbench end to end test", func() {
	Context("creation from samples", func() {
		It("should create sysbench-sample CR", func() {
			_, _, err := run("kubectl create -n " + e2eNamespaceSysbench + " -f " + sysbenchSampleCR)
			Expect(err).To(BeNil())
		})
	})

	Context("created job", func() {
		It("Should finish in a pre-defined time", func() {
			timeout := 120
			cr := &v1alpha1.Sysbench{}
			namespacedName := types.NamespacedName{
				Namespace: e2eNamespaceSysbench,
				Name:      "sysbench-sample",
			}
			Eventually(func() bool {
				if err := client.Get(ctx, namespacedName, cr); err != nil {
					Fail("Unable to get sysbench CR")
				}
				return (cr.Status.Running == false) && (cr.Status.Completed)
			}, timeout).Should(BeTrue())
		})
		It("Should leave a successful job", func() {
			job := &batchv1.Job{}
			namespacedName := types.NamespacedName{
				Namespace: e2eNamespaceSysbench,
				Name:      "sysbench-sample",
			}
			Expect(client.Get(ctx, namespacedName, job)).To(Succeed())
			Expect(job.Status.Succeeded).To(Equal(int32(1)))
		})
	})
})
