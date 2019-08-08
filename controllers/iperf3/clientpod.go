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
	"context"
	"fmt"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/xridge/kubestone/pkg/k8s"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

func (r *Reconciler) newClientPod(ctx context.Context, cr *perfv1alpha1.Iperf3, crRef *corev1.ObjectReference) error {
	iperfCmdLineArgs := fmt.Sprintf("-c %s %s", cr.Name, cr.Spec.ClientConfiguration.CmdLineArgs)
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    cr.Spec.ClientConfiguration.PodLabels,
			Name:      cr.Name + "-client",
			Namespace: cr.Namespace,
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
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
		},
	}

	if err := controllerutil.SetControllerReference(cr, &pod, r.K8S.Scheme); err != nil {
		return err
	}
	if err := r.K8S.Client.Create(ctx, &pod); k8s.IgnoreAlreadyExists(err) != nil {
		return err
	}

	r.K8S.EventRecorder.Eventf(crRef, corev1.EventTypeNormal, k8s.CreateSucceeded,
		"Created Iperf3 Client Pod: %v @ Namespace: %v", pod.Name, pod.Namespace)
	return nil
}
