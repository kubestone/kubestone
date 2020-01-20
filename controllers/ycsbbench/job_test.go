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

package ycsbbench

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("ycsbbench job", func() {
	Describe("NewJob", func() {
		cr := perfv1alpha1.YcsbBench{
			Spec: perfv1alpha1.YcsbBenchSpec{
				Image: perfv1alpha1.ImageSpec{
					Name: "diamantisolutions/ycsb:latest",
				},
				Database: "redis",
				Workload: "a",
				Options:  perfv1alpha1.YcsbBenchOptions{Threadcount: 1, Target: 100},
				Properties: map[string]string{
					"redis.host": "10.0.0.1",
				},
				PodConfig: perfv1alpha1.PodConfigurationSpec{},
			},
		}
		expected_args := []string{
			"redis",
			"-P",
			"workloads/workloada",
			"-threads",
			"1",
			"-target",
			"100",
			"-p",
			"redis.host=10.0.0.1",
		}

		testNewJob(cr, expected_args)
	})

	Describe("NewJob with no target", func() {
		cr := perfv1alpha1.YcsbBench{
			Spec: perfv1alpha1.YcsbBenchSpec{
				Image: perfv1alpha1.ImageSpec{
					Name: "diamantisolutions/ycsb:latest",
				},
				Database: "redis",
				Workload: "a",
				Options:  perfv1alpha1.YcsbBenchOptions{Threadcount: 1},
				Properties: map[string]string{
					"redis.host": "10.0.0.1",
				},
				PodConfig: perfv1alpha1.PodConfigurationSpec{},
			},
		}

		expected_args := []string{
			"redis",
			"-P",
			"workloads/workloada",
			"-threads",
			"1",
			"-p",
			"redis.host=10.0.0.1",
		}

		testNewJob(cr, expected_args)
	})

	Describe("NewJob with no options", func() {
		cr := perfv1alpha1.YcsbBench{
			Spec: perfv1alpha1.YcsbBenchSpec{
				Image: perfv1alpha1.ImageSpec{
					Name: "diamantisolutions/ycsb:latest",
				},
				Database: "redis",
				Workload: "b",
				Properties: map[string]string{
					"redis.host": "10.0.0.1",
				},
				PodConfig: perfv1alpha1.PodConfigurationSpec{},
			},
		}
		expected_args := []string{
			"redis",
			"-P",
			"workloads/workloadb",
			"-p",
			"redis.host=10.0.0.1",
		}

		testNewJob(cr, expected_args)
	})
})

func testNewJob(cr perfv1alpha1.YcsbBench, expected_args []string) {
	job := NewJob(&cr)

	Context("init container", func() {
		containers := job.Spec.Template.Spec.InitContainers
		cont := containers[0]

		It("should have 1 init container", func() {
			Expect(len(containers)).To(Equal(1))
		})
		It("should run load", func() {
			Expect(cont.Command).To(ContainElement("load"))
		})
		Context("arguments", func() {
			It("should format args", func() {
				Expect(cont.Args).To(
					Equal(expected_args))
			})
		})
	})

	Context("main container", func() {
		containers := job.Spec.Template.Spec.Containers
		cont := containers[0]

		It("should have 1 container", func() {
			Expect(len(containers)).To(Equal(1))
		})
		It("should run run", func() {
			Expect(cont.Command).To(ContainElement("run"))
		})
		Context("arguments", func() {
			It("should format args", func() {
				Expect(cont.Args).To(
					Equal(expected_args))
			})
		})
	})
}
