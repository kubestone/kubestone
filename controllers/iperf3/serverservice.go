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

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;create;delete

func newServerService(cr *perfv1alpha1.Iperf3) metav1.Object {
	labels := map[string]string{
		"app":               "iperf3",
		"kubestone-cr-name": cr.Name,
	}
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "iperf3",
					Protocol: "TCP",
					Port:     iperf3ServerPort,
				},
			},
			Selector: labels,
		},
	}

	return &service
}

func (r *Reconciler) deleteServerService(ctx context.Context, cr *perfv1alpha1.Iperf3) error {
	service, err := r.K8S.Clientset.CoreV1().Services(cr.Namespace).Get(cr.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if err := r.K8S.Client.Delete(ctx, service); err != nil {
		return err
	}

	r.K8S.EventRecorder.Eventf(cr, corev1.EventTypeNormal, k8s.DeleteSucceeded,
		"Deleted Iperf3 Server Service: %v @ Namespace: %v", service.Name, service.Namespace)

	return nil
}
