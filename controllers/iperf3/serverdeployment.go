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
	"strconv"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/xridge/kubestone/pkg/k8s"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

const iperf3ServerPort = 5201

func (r *Reconciler) newServerDeployment(ctx context.Context, cr *perfv1alpha1.Iperf3, crRef *corev1.ObjectReference) error {
	replicas := int32(1)
	labels := map[string]string{
		"app":               "iperf3",
		"kubestone-cr-name": cr.Name,
	}
	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "server",
							Image:           "networkstatic/iperf3",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command:         []string{"iperf3"},
							Args:            []string{"-s", "-p", strconv.Itoa(iperf3ServerPort)},
							Ports: []corev1.ContainerPort{
								{
									Name:          "iperf-server",
									ContainerPort: iperf3ServerPort,
								},
							},
							/* -- Causing iperf3 server to exit with 'too many errors'
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.FromInt(iperf3ServerPort),
									},
								},
								InitialDelaySeconds: 5,
								TimeoutSeconds:      5,
								PeriodSeconds:       5,
							},
							*/
						},
					},
				},
			},
		},
	}

	// FIXME: The next three statements are common between serverdeployment, serverservice and clientpod,
	// it would make sense to factor it into one function. For that a type should be found which applies
	// for both SetControllerReference and Create's object.
	if err := controllerutil.SetControllerReference(cr, &deployment, r.K8S.Scheme); err != nil {
		return err
	}
	if err := r.K8S.Client.Create(ctx, &deployment); k8s.IgnoreAlreadyExists(err) != nil {
		return err
	}

	r.K8S.EventRecorder.Eventf(crRef, corev1.EventTypeNormal, k8s.CreateSucceeded,
		"Created Iperf3 Server Deployment: %v @ Namespace: %v", deployment.Name, deployment.Namespace)
	return nil
}

func (r *Reconciler) serverDeploymentReady(cr *perfv1alpha1.Iperf3) (ready bool, err error) {
	ready, err = false, nil
	deployment, err := r.K8S.Clientset.AppsV1().Deployments(cr.Namespace).Get(cr.Name, metav1.GetOptions{})
	if err != nil {
		return ready, err
	}

	if deployment.Status.ReadyReplicas == *deployment.Spec.Replicas {
		ready = true
	}

	return ready, err
}
