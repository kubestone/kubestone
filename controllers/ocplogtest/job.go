package ocplogtest

import (
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
)

// NewJob creates a new ocplogbench job
func NewJob(cr *perfv1alpha1.OcpLogtest) *batchv1.Job {
	objectMeta := metav1.ObjectMeta{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}

	args := []string{
		"ocp_logtest.py",
	}

	if cr.Spec.LineLength > 0 {
		args = append(args, "--line-length", fmt.Sprintf("%d", cr.Spec.LineLength))
	}

	if cr.Spec.NumLines > 0 {
		args = append(args, "--num-lines", fmt.Sprintf("%d", cr.Spec.NumLines))
	}

	if cr.Spec.Rate > 0 {
		args = append(args, "--rate", fmt.Sprintf("%d", cr.Spec.Rate))
	}

	if cr.Spec.FixedLine {
		args = append(args, "--fixed-line")
	}

	job := k8s.NewPerfJob(objectMeta, "ocplogbench", cr.Spec.Image, cr.Spec.PodConfig)
	job.Spec.Template.Spec.Containers[0].Command = []string{"python"}
	job.Spec.Template.Spec.Containers[0].Args = args

	return job
}
