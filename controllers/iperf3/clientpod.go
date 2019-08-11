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
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;create;delete

func clientPodName(cr *perfv1alpha1.Iperf3) string {
	// Should not match with service name as the pod's
	// hostname is set to it's name. If the two matches
	// the destination ip will resolve to 127.0.0.1 and
	// the server will be unreachable.
	return serverServiceName(cr) + "-client"
}

// NewClientPod creates an Iperf3 Client Pod (targetting the
// Server Deployment via the Server Service) from the provided
// IPerf3 Benchmark Definition.
func NewClientPod(cr *perfv1alpha1.Iperf3) *corev1.Pod {
	serverAddress := serverServiceName(cr)
	iperfCmdLineArgs := fmt.Sprintf("--client %s --port %d ",
		serverAddress, Iperf3ServerPort)

	if cr.Spec.UDP {
		iperfCmdLineArgs += "--udp "
	}

	iperfCmdLineArgs += cr.Spec.ClientConfiguration.CmdLineArgs
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    cr.Spec.ClientConfiguration.PodLabels,
			Name:      clientPodName(cr),
			Namespace: cr.Namespace,
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			ImagePullSecrets: []corev1.LocalObjectReference{
				{
					Name: cr.Spec.Image.PullSecret,
				},
			},
			Containers: []corev1.Container{
				{
					Name:            "client",
					Image:           cr.Spec.Image.Name,
					ImagePullPolicy: corev1.PullPolicy(cr.Spec.Image.PullPolicy),
					Command:         []string{"iperf3"},
					Args:            strings.Fields(iperfCmdLineArgs),
				},
			},
			Affinity:     &cr.Spec.ClientConfiguration.PodScheduling.Affinity,
			Tolerations:  cr.Spec.ClientConfiguration.PodScheduling.Tolerations,
			NodeSelector: cr.Spec.ClientConfiguration.PodScheduling.NodeSelector,
			NodeName:     cr.Spec.ClientConfiguration.PodScheduling.NodeName,
			HostNetwork:  cr.Spec.ClientConfiguration.HostNetwork,
		},
	}

	return &pod
}

func (r *Reconciler) clientPodFinished(cr *perfv1alpha1.Iperf3) (finished bool, err error) {
	finished, err = false, nil
	pod, err := r.K8S.Clientset.CoreV1().Pods(cr.Namespace).Get(clientPodName(cr), metav1.GetOptions{})
	if err != nil {
		return finished, err
	}

	finished = pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed

	return finished, err
}
