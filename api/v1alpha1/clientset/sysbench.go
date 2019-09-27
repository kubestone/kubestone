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

// SysbenchesGetter has a method to return a SysbenchInterface.
// A group's client should implement this interface.
type SysbenchesGetter interface {
	Sysbenches(namespace string) SysbenchInterface
}

// SysbenchInterface has methods to work with Sysbench resources.
type SysbenchInterface interface {
	Create(*perfv1alpha1.Sysbench) (*perfv1alpha1.Sysbench, error)
	Update(*perfv1alpha1.Sysbench) (*perfv1alpha1.Sysbench, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*perfv1alpha1.Sysbench, error)
	List(opts metav1.ListOptions) (*perfv1alpha1.SysbenchList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *perfv1alpha1.Sysbench, err error)
	SysbenchExpansion
}

// sysbenches implements SysbenchInterface
type sysbenches struct {
	client rest.Interface
	ns     string
}

// newSysbenches returns a Sysbenches
func newSysbenches(c *PerfV1alpha1Client, namespace string) *sysbenches {
	return &sysbenches{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the sysbench, and returns the corresponding sysbench object, and an error if there is any.
func (c *sysbenches) Get(name string, options metav1.GetOptions) (result *perfv1alpha1.Sysbench, err error) {
	result = &perfv1alpha1.Sysbench{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sysbenches").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Sysbenches that match those selectors.
func (c *sysbenches) List(opts metav1.ListOptions) (result *perfv1alpha1.SysbenchList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &perfv1alpha1.SysbenchList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sysbenches").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested sysbenches.
func (c *sysbenches) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("sysbenches").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a sysbench and creates it.  Returns the server's representation of the sysbench, and an error, if there is any.
func (c *sysbenches) Create(sysbench *perfv1alpha1.Sysbench) (result *perfv1alpha1.Sysbench, err error) {
	result = &perfv1alpha1.Sysbench{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("sysbenches").
		Body(sysbench).
		Do().
		Into(result)
	return
}

// Update takes the representation of a sysbench and updates it. Returns the server's representation of the sysbench, and an error, if there is any.
func (c *sysbenches) Update(sysbench *perfv1alpha1.Sysbench) (result *perfv1alpha1.Sysbench, err error) {
	result = &perfv1alpha1.Sysbench{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("sysbenches").
		Name(sysbench.Name).
		Body(sysbench).
		Do().
		Into(result)
	return
}

// Delete takes name of the sysbench and deletes it. Returns an error if one occurs.
func (c *sysbenches) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sysbenches").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *sysbenches) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sysbenches").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched sysbench.
func (c *sysbenches) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *perfv1alpha1.Sysbench, err error) {
	result = &perfv1alpha1.Sysbench{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("sysbenches").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
