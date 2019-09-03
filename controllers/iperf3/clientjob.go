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

package iperf3

import (
	"strconv"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/firepear/qsplit"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;create;delete

func clientJobName(cr *perfv1alpha1.Iperf3) string {
	// Should not match with service name as the pod's
	// hostname is set to it's name. If the two matches
	// the destination ip will resolve to 127.0.0.1 and
	// the server will be unreachable.
	return serverServiceName(cr) + "-client"
}

// NewClientJob creates an Iperf3 Client Job (targeting the
// Server Deployment via the Server Service) from the provided
// IPerf3 Benchmark Definition.
func NewClientJob(cr *perfv1alpha1.Iperf3) *batchv1.Job {
	serverAddress := serverServiceName(cr)
	iperfCmdLineArgs := []string{
		"--client", serverAddress,
		"--port", strconv.Itoa(Iperf3ServerPort),
	}

	if cr.Spec.UDP {
		iperfCmdLineArgs = append(iperfCmdLineArgs, "--udp")
	}

	iperfCmdLineArgs = append(iperfCmdLineArgs,
		qsplit.ToStrings([]byte(cr.Spec.ClientConfiguration.CmdLineArgs))...)

	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clientJobName(cr),
			Namespace: cr.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: cr.Spec.ClientConfiguration.PodLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "client",
							Image:           cr.Spec.Image.Name,
							ImagePullPolicy: corev1.PullPolicy(cr.Spec.Image.PullPolicy),
							Command:         []string{"iperf3"},
							Args:            iperfCmdLineArgs,
						},
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: cr.Spec.Image.PullSecret,
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Affinity:      &cr.Spec.ClientConfiguration.PodScheduling.Affinity,
					Tolerations:   cr.Spec.ClientConfiguration.PodScheduling.Tolerations,
					NodeSelector:  cr.Spec.ClientConfiguration.PodScheduling.NodeSelector,
					NodeName:      cr.Spec.ClientConfiguration.PodScheduling.NodeName,
					HostNetwork:   cr.Spec.ClientConfiguration.HostNetwork,
				},
			},
		},
	}

	return &job
}
