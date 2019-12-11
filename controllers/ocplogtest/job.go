package ocplogtest

import (
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
)

// NewJob creates a new pgbench job
func NewJob(cr *perfv1alpha1.OcpLogtest) *batchv1.Job {
	objectMeta := metav1.ObjectMeta{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}

	job := k8s.NewPerfJob(objectMeta, "pgbench", cr.Spec.Image, cr.Spec.PodConfig)

	args := []string{
		"ocp_logtest.py",
	}

	if cr.Spec.FixedLine == true {
		args = append(args, "--line-length")
	}

	if cr.Spec.LineLength > 0 {
		args = append(args, fmt.Sprintf("--fixed-line=%d", cr.Spec.LineLength))
	}

	if cr.Spec.NumLines > 0 {
		args = append(args, fmt.Sprintf("--num-lines=%d", cr.Spec.NumLines))
	}

	if cr.Spec.Rate > 0 {
		args = append(args, fmt.Sprintf("--rate=%d", cr.Spec.Rate))
	}

	job.Spec.Template.Spec.Containers[0].Command = []string{"python"}
	job.Spec.Template.Spec.Containers[0].Args = args

	return job
}
