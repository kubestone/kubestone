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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	drillSampleCR = "../../config/samples/perf_v1alpha1_drill.yaml"
)

var _ = Describe("end to end test", func() {
	Context("create from samples", func() {
		It("should create drill-sample cr", func() {
			_, _, err := run("kubectl create -n " + e2eNamespaceDrill + " -f " + drillSampleCR)
			Expect(err).To(BeNil())
		})
	})

	Context("created job", func() {
		It("Should finish in a pre-defined time", func() {
			timeout := 60
			cr := &v1alpha1.Drill{}
			namespacedName := types.NamespacedName{
				Namespace: e2eNamespaceDrill,
				Name:      "drill-sample",
			}
			Eventually(func() bool {
				if err := client.Get(ctx, namespacedName, cr); err != nil {
					Fail("Unable to get drill CR")
				}
				return (cr.Status.Running == false) && (cr.Status.Completed)
			}, timeout).Should(BeTrue())
		})
		It("Should leave one successful pod that actually fetched kubernetes.io", func() {
			inNamespace := ctrlclient.InNamespace(e2eNamespaceDrill)
			matchingLabel := ctrlclient.MatchingLabels{
				"job-name": "drill-sample",
			}
			podList := &corev1.PodList{}
			Expect(client.List(ctx, podList, inNamespace, matchingLabel)).To(Succeed())
			Expect(len(podList.Items)).To(Equal(1))
			pod := podList.Items[0]

			output, _, err := run("kubectl logs -n " + e2eNamespaceDrill + " " + pod.Name)
			Expect(err).To(BeNil())
			Expect(output).To(ContainSubstring("Fetch kubernetes.io"))
		})
		It("Should leave one successful job", func() {
			job := &batchv1.Job{}
			namespacedName := types.NamespacedName{
				Namespace: e2eNamespaceDrill,
				Name:      "drill-sample",
			}
			Expect(client.Get(ctx, namespacedName, job)).To(Succeed())
			Expect(job.Status.Succeeded).To(Equal(int32(1)))
		})
	})
})
