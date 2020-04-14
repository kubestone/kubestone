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

package esrally

import (
	"fmt"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func NewJob(cr *perfv1alpha1.EsRally) *batchv1.Job {
	objectMeta := metav1.ObjectMeta{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}

	job := k8s.NewPerfJob(objectMeta, "esrally", cr.Spec.Image, cr.Spec.PodConfig)
	job.Spec.Template.Spec.Containers[0].Command = []string{"/bin/sh", "-c"}
	job.Spec.Template.Spec.Containers[0].Args = []string{
		"touch /rally/.rally/logs/rally.log && tail -f /rally/.rally/logs/rally.log & " +
			strings.Join(CreateEsRallyCmd(&cr.Spec, &objectMeta), " "),
	}

	return job
}

func CreateEsRallyCmd(spec *perfv1alpha1.EsRallySpec, objectMeta *metav1.ObjectMeta) []string {
	var cmdArgs []string

	cmdArgs = append(cmdArgs, "/usr/local/bin/esrally", "--pipeline=benchmark-only")
	cmdArgs = append(cmdArgs, fmt.Sprintf("--track=%s", spec.Track))
	cmdArgs = append(cmdArgs, fmt.Sprintf("--challenge=%s", spec.Challenge))
	cmdArgs = append(cmdArgs, fmt.Sprintf("--target-hosts=%s", spec.Hosts))
	cmdArgs = append(cmdArgs, fmt.Sprintf("--load-driver-hosts=%s", ParseRallyNodeNames(spec, objectMeta)))

	if spec.TrackRepository != nil {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--track-repository=%s", *spec.TrackRepository))
	}

	if spec.TrackParams != nil {
		var params string
		for key, val := range *spec.TrackParams {
			params = params + fmt.Sprintf("%s:%s,", key, val)
		}
		cmdArgs = append(cmdArgs, fmt.Sprintf("--track-params=%s", strings.Trim(params, ",")))
	}

	return cmdArgs
}

func ParseRallyNodeNames(spec *perfv1alpha1.EsRallySpec, objectMeta *metav1.ObjectMeta) string {
	var nodes string
	var nodeCount int32
	if spec.Nodes == nil {
		nodeCount = 1
	} else {
		nodeCount = *spec.Nodes
	}

	for i := int32(0); i < nodeCount; i++ {
		nodes = nodes + fmt.Sprintf("%s-%d.%s,", objectMeta.Name, i, objectMeta.Name)
	}

	return strings.Trim(nodes, ",")
}
