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

package kafkabench

import (
	"fmt"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func NewConsumerJob(cr *perfv1alpha1.KafkaBench, ts *perfv1alpha1.KafkaTestSpec) *batchv1.Job {
	jobName := fmt.Sprintf("%s-%s-consumer", cr.Name, ts.Name)

	objectMeta := metav1.ObjectMeta{
		Name:      jobName,
		Namespace: cr.Namespace,
	}

	job := k8s.NewPerfJob(objectMeta, "kafkabench", cr.Spec.Image, cr.Spec.PodConfig)
	job.Spec.Parallelism = &ts.Threads

	consumerSleep := int32(40)

	if ts.ConsumerSleep != nil {
		consumerSleep = *ts.ConsumerSleep
	}

	// Add init job to sleep, this allows the producer to queue up messages
	initContainer := corev1.Container{
		Name:            "kafka-consumer-init",
		Image:           cr.Spec.Image.Name,
		ImagePullPolicy: corev1.PullPolicy(cr.Spec.Image.PullPolicy),
		Command:         []string{"/bin/sleep", fmt.Sprintf("%d", consumerSleep)},
		Resources:       cr.Spec.PodConfig.Resources,
	}
	job.Spec.Template.Spec.InitContainers = append(job.Spec.Template.Spec.InitContainers, initContainer)

	job.Spec.Template.Spec.Containers[0].Command = []string{"/bin/sh"}
	job.Spec.Template.Spec.Containers[0].Args = ConsumerJobArgs(cr, ts)

	// Add pod affinity
	AddPodAffinity(job, jobName)

	return job
}

func ConsumerJobArgs(cr *perfv1alpha1.KafkaBench, ts *perfv1alpha1.KafkaTestSpec) []string {
	brokers := strings.Join(cr.Spec.Brokers, ",")

	timeout := "10000"
	if ts.Timeout != nil {
		timeout = fmt.Sprintf("%d", *ts.Timeout)
	}

	return []string{
		"/usr/bin/kafka-consumer-perf-test",
		"--broker-list", brokers,
		"--messages", fmt.Sprintf("%d", ts.Records),
		"--threads", "1",
		"--topic", fmt.Sprintf("%s-%s-bench", cr.Name, ts.Name),
		"--timeout", timeout,
	}
}
