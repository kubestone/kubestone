package ycsbbench

import (
	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/pkg/k8s"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewJob(cr *perfv1alpha1.YcsbBench) *batchv1.Job {
	objectMeta := metav1.ObjectMeta{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}

	env := []corev1.EnvVar{
		{Name: "ACTION", Value: cr.Spec.Action},
		{Name: "DBTYPE", Value: cr.Spec.DbType},
		{Name: "WORKLETTER", Value: cr.Spec.Workletter},
		{Name: "DBARGS", Value: cr.Spec.DbArgs},
	}

	job := k8s.NewPerfJob(objectMeta, "ycsbbench", cr.Spec.Image, cr.Spec.PodConfig)
	job.Spec.Template.Spec.Containers[0].Env = env

	return job
}
