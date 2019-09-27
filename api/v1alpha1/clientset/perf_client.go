package clientset

import (
	"k8s.io/apimachinery/pkg/runtime/serializer"
	scheme "k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

func init() {
	_ = perfv1alpha1.AddToScheme(scheme.Scheme)
}

type PerfV1alpha1Interface interface {
	RESTClient() rest.Interface
	DrillsGetter
	FiosGetter
	IopingsGetter
	Iperf3sGetter
	PgbenchesGetter
	QpervesGetter
	SysbenchesGetter
}

// PerfV1alpha1Client is used to interact with features provided by the perf.kubestone.xridge.io group.
type PerfV1alpha1Client struct {
	restClient rest.Interface
}

func (c *PerfV1alpha1Client) Drills(namespace string) DrillInterface {
	return newDrills(c, namespace)
}

func (c *PerfV1alpha1Client) Fios(namespace string) FioInterface {
	return newFios(c, namespace)
}

func (c *PerfV1alpha1Client) Iopings(namespace string) IopingInterface {
	return newIopings(c, namespace)
}

func (c *PerfV1alpha1Client) Iperf3s(namespace string) Iperf3Interface {
	return newIperf3s(c, namespace)
}

func (c *PerfV1alpha1Client) Pgbenches(namespace string) PgbenchInterface {
	return newPgbenches(c, namespace)
}

func (c *PerfV1alpha1Client) Qperves(namespace string) QperfInterface {
	return newQperves(c, namespace)
}

func (c *PerfV1alpha1Client) Sysbenches(namespace string) SysbenchInterface {
	return newSysbenches(c, namespace)
}

// NewForConfig creates a new PerfV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*PerfV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &PerfV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new PerfV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *PerfV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

func setConfigDefaults(config *rest.Config) error {
	gv := perfv1alpha1.GroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *PerfV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
