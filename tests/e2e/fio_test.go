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
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/xridge/kubestone/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	fioCrBaseDir = samplesDir + "/fio"
)

var _ = FDescribe("end to end test", func() {
	DescribeTable("creating fio job from multiple CRs",
		func(crDir string) {
			splits := strings.Split(crDir, "/")
			dirName := splits[len(splits)-1]
			crName := "fio-" + strings.ReplaceAll(dirName, "_", "-")

			By("creating fio job from " + crDir)
			_, _, err := run(`bash -c "` +
				"kustomize build " + crDir + " | " +
				"sed 's/name: fio-sample/name: " + crName + "/' | " +
				"kubectl create -n " + e2eNamespaceFio + ` -f -"`)
			Expect(err).To(BeNil())

			By("checking the created CR")
			timeout := 180
			cr := &v1alpha1.Fio{}
			namespacedName := types.NamespacedName{
				Namespace: e2eNamespaceFio,
				Name:      crName,
			}
			Eventually(func() bool {
				if err := client.Get(ctx, namespacedName, cr); err != nil {
					Fail("Unable to get fio CR: " + err.Error())
				}
				return !cr.Status.Running && cr.Status.Completed
			}, timeout).Should(BeTrue())

			By("checking the created job")
			job := &batchv1.Job{}
			Expect(client.Get(ctx, namespacedName, job)).To(Succeed())
			Expect(job.Status.Succeeded).To(Equal(int32(1)))
		},

		Entry("base", fioCrBaseDir+"/base"),
		Entry("emptydir", fioCrBaseDir+"/overlays/emptydir"),
		Entry("builtin jobs", fioCrBaseDir+"/overlays/builtin_jobs"),
		Entry("custom jobs", fioCrBaseDir+"/overlays/custom_jobs"),
	)
})
