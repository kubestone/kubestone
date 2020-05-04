package kafkabench

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ksapi "github.com/xridge/kubestone/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

var _ = Describe("Consumer Job", func() {
	Describe("created from CR", func() {
		var cr ksapi.KafkaBench
		var jobs []*batchv1.Job
		BeforeEach(func() {
			cr = ksapi.KafkaBench{
				ObjectMeta: v1.ObjectMeta{
					Name: "kafkabench-sample",
				},
				Spec: ksapi.KafkaBenchSpec{
					Image: ksapi.ImageSpec{
						Name: "confluentinc/cp-kafka:5.2.1",
					},
					PodConfig: ksapi.PodConfigurationSpec{},
					KafkaClusterInfo: ksapi.KafkaClusterInfo{
						ZooKeepers: []string{
							"kafka-demo-cp-zookeeper-0.kafka-demo-cp-zookeeper-headless.kafka-demo.svc.cluster.local:2181",
							"kafka-demo-cp-zookeeper-1.kafka-demo-cp-zookeeper-headless.kafka-demo.svc.cluster.local:2181",
							"kafka-demo-cp-zookeeper-2.kafka-demo-cp-zookeeper-headless.kafka-demo.svc.cluster.local:2181",
						},
						Brokers: []string{
							"kafka-demo-cp-kafka-0.kafka-demo-cp-kafka-headless.kafka-demo.svc.cluster.local:9092",
							"kafka-demo-cp-kafka-1.kafka-demo-cp-kafka-headless.kafka-demo.svc.cluster.local:9092",
							"kafka-demo-cp-kafka-2.kafka-demo-cp-kafka-headless.kafka-demo.svc.cluster.local:9092",
						},
					},
					Tests: []ksapi.KafkaTestSpec{
						{
							Name:        "noreplication",
							Threads:     2,
							Replication: 1,
							Partitions:  16,
							RecordSize:  100,
							Records:     60000000,
							ExtraProducerOpts: []string{
								"buffer.memory=671088640",
							},
							ConsumersOnly: false,
							ProducersOnly: false,
						},
						{
							Name:        "replication",
							Threads:     3,
							Replication: 3,
							Partitions:  16,
							RecordSize:  100,
							Records:     60000000,
							ExtraProducerOpts: []string{
								"acks=1",
								"buffer.memory=671088640",
							},
							ConsumersOnly: false,
							ProducersOnly: false,
						},
					},
				},
			}

			for _, test := range cr.Spec.Tests {
				jobs = append(jobs, NewConsumerJob(&cr, &test))
			}
		})

		Context("with multiple tests", func() {
			It("create multiple jobs", func() {
				Expect(jobs).To(HaveLen(2))
			})
		})

		Context("with ClusterInfo specified", func() {
			It("should contain broker list", func() {
				Expect(jobs[0].Spec.Template.Spec.Containers[0].Args).To(
					ContainElement(strings.Join(cr.Spec.KafkaClusterInfo.Brokers, ",")))
			})
		})
	})
})
