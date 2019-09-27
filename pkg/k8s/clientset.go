package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/xridge/kubestone/api/v1alpha1/clientset"
)

type Interface interface {
	PerfV1alpha1() clientset.PerfV1alpha1Interface
}

type Clientset struct {
	kubernetes.Clientset
	perfv1alpha1 *clientset.PerfV1alpha1Client
}

var _ Interface = &Clientset{}

// PerfV1alpha1 retrieves the PerfV1alpha1Client
func (c *Clientset) PerfV1alpha1() clientset.PerfV1alpha1Interface {
	return c.perfv1alpha1
}

func NewForConfig(config *rest.Config) (*Clientset, error) {
	kcs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	pcs, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Clientset{*kcs, pcs}, nil
}
