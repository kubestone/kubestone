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
	pgBenchSampleCR = samplesDir + "/perf_v1alpha1_pgbench.yaml"
)

var _ = Describe("pgbench end to end test", func() {
	Context("install postgresql", func() {
		It("should succeed", func() {
			_, _, err := run("kubectl create -f " + testConf + "/postgres.yaml -n " + e2eNamespacePgbench)
			Expect(err).To(BeNil())

			By("wait until postgres actually starts")
			_, _, err = run("kubectl rollout status statefulset/postgres -n " + e2eNamespacePgbench)
			Expect(err).To(BeNil())
		})
	})

	Context("create from the example", func() {
		It("should succeed", func() {
			_, _, err := run("kubectl create -f " + pgBenchSampleCR + " -n " + e2eNamespacePgbench)
			Expect(err).To(BeNil())
		})
	})

	Context("created job", func() {
		It("should finish in a pre-defined time", func() {
			timeout := 60
			cr := &v1alpha1.Pgbench{}
			namespacedName := types.NamespacedName{
				Namespace: e2eNamespacePgbench,
				Name:      "pgbench-sample",
			}
			Eventually(func() bool {
				if err := client.Get(ctx, namespacedName, cr); err != nil {
					Fail("Unable to get pgbench CR")
				}
				return (cr.Status.Running == false) && (cr.Status.Completed)
			}, timeout).Should(BeTrue())
		})

		It("should leave a successful job", func() {
			pod := &batchv1.Job{}
			namespacedName := types.NamespacedName{
				Namespace: e2eNamespacePgbench,
				Name:      "pgbench-sample",
			}
			Expect(client.Get(ctx, namespacedName, pod)).To(Succeed())
			Expect(pod.Status.Succeeded).To(Equal(int32(1)))
		})
	})
})
