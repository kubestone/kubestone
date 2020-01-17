package ycsbbench

import (
	"fmt"
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

func formatArgs(cr *perfv1alpha1.YcsbBench) []string {
	args := []string{
		cr.Spec.Database,
		"-P", fmt.Sprintf("workloads/workload%s", cr.Spec.Workload),
	}

	if cr.Spec.Options != (perfv1alpha1.YcsbBenchOptions{}) {
		if cr.Spec.Options.Threadcount > 0 {
			args = append(args, "-threads", strconv.Itoa(cr.Spec.Options.Threadcount))
		}

		if cr.Spec.Options.Target > 0 {
			args = append(args, "-target", strconv.Itoa(cr.Spec.Options.Target))
		}
	}

	for key, val := range cr.Spec.Properties {
		args = append(args, "-p", fmt.Sprintf("%s=%s", key, val))
	}
	return args
}

func NewJob(cr *perfv1alpha1.YcsbBench) *batchv1.Job {
	objectMeta := metav1.ObjectMeta{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}

	args := formatArgs(cr)
	initContainer := corev1.Container{
		Name:            "ycsbbench-load",
		Image:           cr.Spec.Image.Name,
		ImagePullPolicy: corev1.PullPolicy(cr.Spec.Image.PullPolicy),
		Command:         []string{"./bin/ycsb", "load"},
		Args:            args,
	}
	// append([]string{"./bin/ycsb", "load", args},
	job := k8s.NewPerfJob(objectMeta, "ycsbbench", cr.Spec.Image, cr.Spec.PodConfig)
	job.Spec.Template.Spec.InitContainers = append(
		job.Spec.Template.Spec.InitContainers, initContainer)

	job.Spec.Template.Spec.Containers[0].Command = []string{"./bin/ycsb", "run"}
	job.Spec.Template.Spec.Containers[0].Args = args

	return job
}
