package clientset

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	scheme "k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

// Iperf3sGetter has a method to return a Iperf3Interface.
// A group's client should implement this interface.
type Iperf3sGetter interface {
	Iperf3s(namespace string) Iperf3Interface
}

// Iperf3Interface has methods to work with Iperf3 resources.
type Iperf3Interface interface {
	Create(*perfv1alpha1.Iperf3) (*perfv1alpha1.Iperf3, error)
	Update(*perfv1alpha1.Iperf3) (*perfv1alpha1.Iperf3, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*perfv1alpha1.Iperf3, error)
	List(opts metav1.ListOptions) (*perfv1alpha1.Iperf3List, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *perfv1alpha1.Iperf3, err error)
	Iperf3Expansion
}

// iperf3s implements Iperf3Interface
type iperf3s struct {
	client rest.Interface
	ns     string
}

// newIperf3s returns a Iperf3s
func newIperf3s(c *PerfV1alpha1Client, namespace string) *iperf3s {
	return &iperf3s{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the iperf3, and returns the corresponding iperf3 object, and an error if there is any.
func (c *iperf3s) Get(name string, options metav1.GetOptions) (result *perfv1alpha1.Iperf3, err error) {
	result = &perfv1alpha1.Iperf3{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("iperf3s").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Iperf3s that match those selectors.
func (c *iperf3s) List(opts metav1.ListOptions) (result *perfv1alpha1.Iperf3List, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &perfv1alpha1.Iperf3List{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("iperf3s").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested iperf3s.
func (c *iperf3s) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("iperf3s").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a iperf3 and creates it.  Returns the server's representation of the iperf3, and an error, if there is any.
func (c *iperf3s) Create(iperf3 *perfv1alpha1.Iperf3) (result *perfv1alpha1.Iperf3, err error) {
	result = &perfv1alpha1.Iperf3{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("iperf3s").
		Body(iperf3).
		Do().
		Into(result)
	return
}

// Update takes the representation of a iperf3 and updates it. Returns the server's representation of the iperf3, and an error, if there is any.
func (c *iperf3s) Update(iperf3 *perfv1alpha1.Iperf3) (result *perfv1alpha1.Iperf3, err error) {
	result = &perfv1alpha1.Iperf3{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("iperf3s").
		Name(iperf3.Name).
		Body(iperf3).
		Do().
		Into(result)
	return
}

// Delete takes name of the iperf3 and deletes it. Returns an error if one occurs.
func (c *iperf3s) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("iperf3s").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *iperf3s) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("iperf3s").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched iperf3.
func (c *iperf3s) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *perfv1alpha1.Iperf3, err error) {
	result = &perfv1alpha1.Iperf3{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("iperf3s").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
