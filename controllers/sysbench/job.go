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

package sysbench

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/firepear/qsplit"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
)

// NewJob creates a sysbench benchmark job
func NewJob(cr *perfv1alpha1.Sysbench) *batchv1.Job {
	objectMeta := metav1.ObjectMeta{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}

	sysbenchCmdLineArgs := []string{}
	sysbenchCmdLineArgs = append(sysbenchCmdLineArgs, qsplit.ToStrings([]byte(cr.Spec.Options))...)
	sysbenchCmdLineArgs = append(sysbenchCmdLineArgs, cr.Spec.TestName, cr.Spec.Command)

	job := k8s.NewPerfJob(objectMeta, "sysbench", cr.Spec.Image, cr.Spec.PodConfig)
	job.Spec.Template.Spec.Containers[0].Args = sysbenchCmdLineArgs
	return job
}
