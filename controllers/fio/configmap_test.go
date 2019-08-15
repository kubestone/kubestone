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
package fio

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

const jobFile0 = `
[global]
ioengine=rbd
clientname=admin
pool=rbd
rbdname=fio_test
rw=randwrite
bs=4k

[rbd_iodepth32]
iodepth=32`

const jobFile1 = `
[random-writers]
ioengine=libaio
iodepth=4
rw=randwrite
bs=32k
direct=0
size=64m
numjobs=4`

var _ = Describe("fio configmap", func() {
	Describe("cr with custom job files", func() {
		var cr perfv1alpha1.Fio
		var configMap *corev1.ConfigMap

		BeforeEach(func() {
			cr = perfv1alpha1.Fio{
				Spec: perfv1alpha1.FioSpec{
					CustomJobFiles: []string{jobFile0, jobFile1},
				},
			}
			configMap = NewConfigMap(&cr)
		})

		Context("with custom job files", func() {
			It("should have them in the data", func() {
				Expect(configMap.Data[CustomJobName(0)]).To(
					Equal(jobFile0))
				Expect(configMap.Data[CustomJobName(1)]).To(
					Equal(jobFile1))
			})
		})
	})
})
