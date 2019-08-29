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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	iperf3SampleCR     = "../../config/samples/perf_v1alpha1_iperf3.yaml"
	e2eNamespaceIperf3 = "kubestone-e2e-iperf3"
)

var _ = Describe("end to end test", func() {
	Context("preparing namespace", func() {
		_, _, err := run("kubectl create namespace " + e2eNamespaceIperf3)
		It("should succeed", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("for iperf3", func() {
		Context("creation from samples", func() {
			_, _, err := run("kubectl create -n " + e2eNamespaceIperf3 + " -f " + iperf3SampleCR)
			It("should create iperf3-sample cr", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("created job", func() {
			It("Should finish in a pre-defined time", func() {
				timeout := 60
				cr := &v1alpha1.Iperf3{}
				// TODO: find the respective objects via the CR owner reference
				namespacedName := types.NamespacedName{
					Namespace: e2eNamespaceIperf3,
					Name:      "iperf3-sample",
				}
				Eventually(func() bool {
					if err := client.Get(ctx, namespacedName, cr); err != nil {
						Fail("Unable to get iperf3 CR")
					}
					return (cr.Status.Running == false) && (cr.Status.Completed)
				}, timeout).Should(BeTrue())
			})
			It("Should leave one successful pod", func() {
				pod := &corev1.Pod{}
				namespacedName := types.NamespacedName{
					Namespace: e2eNamespaceIperf3,
					Name:      "iperf3-sample-client",
				}
				Expect(client.Get(ctx, namespacedName, pod)).To(Succeed())
				Expect(pod.Status.Phase).To(Equal(corev1.PodSucceeded))
			})
			It("Should not leave deployment", func() {
				deployment := &appsv1.Deployment{}
				namespacedName := types.NamespacedName{
					Namespace: e2eNamespaceIperf3,
					Name:      "iperf3-sample",
				}
				Expect(client.Get(ctx, namespacedName, deployment)).NotTo(Succeed())
			})
			It("Should not leave service", func() {
				service := &corev1.Service{}
				namespacedName := types.NamespacedName{
					Namespace: e2eNamespaceIperf3,
					Name:      "iperf3-sample",
				}
				Expect(client.Get(ctx, namespacedName, service)).NotTo(Succeed())
			})
		})
	})
})
