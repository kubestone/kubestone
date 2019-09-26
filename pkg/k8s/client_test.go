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

package k8s

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("Client", func() {

	const namespace = "fake-namespace"

	var access *Access
	var pod, badPod *corev1.Pod
	var fioCr *perfv1alpha1.Fio

	BeforeEach(func(done Done) {
		cl, err := client.New(cfg, client.Options{})
		Expect(err).NotTo(HaveOccurred())
		Expect(cl).NotTo(BeNil())

		access = &Access{
			Client:        cl,
			Clientset:     clientset,
			Scheme:        scheme,
			EventRecorder: NewEventRecorder(clientset, nil),
		}
		pod = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fake-pod",
				Namespace: namespace,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "fake-container",
						Image: "fake-image",
					},
				},
			},
		}
		badPod = &corev1.Pod{ // has missing required fields
			ObjectMeta: metav1.ObjectMeta{Name: "fake-pod", Namespace: namespace},
			Spec:       corev1.PodSpec{},
		}
		fioCr = &perfv1alpha1.Fio{
			TypeMeta: metav1.TypeMeta{
				Kind:       "fio",
				APIVersion: "v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fake-fio",
				Namespace: namespace,
				UID:       types.UID("cr-uid"),
			},
		}
		close(done)
	}, 10)

	Describe("CreateWithReference", func() {

		Context("with valid object and owner", func() {
			It("should create a new object and set the reference", func() {
				err := access.CreateWithReference(context.TODO(), pod, fioCr)
				Expect(err).NotTo(HaveOccurred())

				By("check the created pod")
				actual, err := clientset.CoreV1().Pods(pod.Namespace).Get(pod.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(actual).NotTo(BeNil())
				Expect(pod).To(Equal(actual))

				By("check the owner reference")
				ownerReferences := pod.GetOwnerReferences()
				Expect(ownerReferences).NotTo(BeEmpty())

				By("check the created event")
				Eventually(func() []corev1.Event {
					eventList, err := clientset.CoreV1().Events(pod.Namespace).List(metav1.ListOptions{})
					Expect(err).NotTo(HaveOccurred())
					Expect(eventList).NotTo(BeNil())

					if len(eventList.Items) > 0 {
						event := eventList.Items[0]
						Expect(event.Reason).To(Equal(Created))
						Expect(event.Message).To(ContainSubstring("Created"))
						Expect(event.Message).To(ContainSubstring(pod.Name))
					}

					return eventList.Items
				}).ShouldNot(BeEmpty())
			})
		})

		Context("with invalid object", func() {
			It("should fail", func() {
				err := access.CreateWithReference(context.TODO(), badPod, fioCr)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
