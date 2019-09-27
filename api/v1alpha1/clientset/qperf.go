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

// QpervesGetter has a method to return a QperfInterface.
// A group's client should implement this interface.
type QpervesGetter interface {
	Qperves(namespace string) QperfInterface
}

// QperfInterface has methods to work with Qperf resources.
type QperfInterface interface {
	Create(*perfv1alpha1.Qperf) (*perfv1alpha1.Qperf, error)
	Update(*perfv1alpha1.Qperf) (*perfv1alpha1.Qperf, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*perfv1alpha1.Qperf, error)
	List(opts metav1.ListOptions) (*perfv1alpha1.QperfList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *perfv1alpha1.Qperf, err error)
	QperfExpansion
}

// qperves implements QperfInterface
type qperves struct {
	client rest.Interface
	ns     string
}

// newQperves returns a Qperves
func newQperves(c *PerfV1alpha1Client, namespace string) *qperves {
	return &qperves{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the qperf, and returns the corresponding qperf object, and an error if there is any.
func (c *qperves) Get(name string, options metav1.GetOptions) (result *perfv1alpha1.Qperf, err error) {
	result = &perfv1alpha1.Qperf{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("qperves").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Qperves that match those selectors.
func (c *qperves) List(opts metav1.ListOptions) (result *perfv1alpha1.QperfList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &perfv1alpha1.QperfList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("qperves").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested qperves.
func (c *qperves) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("qperves").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a qperf and creates it.  Returns the server's representation of the qperf, and an error, if there is any.
func (c *qperves) Create(qperf *perfv1alpha1.Qperf) (result *perfv1alpha1.Qperf, err error) {
	result = &perfv1alpha1.Qperf{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("qperves").
		Body(qperf).
		Do().
		Into(result)
	return
}

// Update takes the representation of a qperf and updates it. Returns the server's representation of the qperf, and an error, if there is any.
func (c *qperves) Update(qperf *perfv1alpha1.Qperf) (result *perfv1alpha1.Qperf, err error) {
	result = &perfv1alpha1.Qperf{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("qperves").
		Name(qperf.Name).
		Body(qperf).
		Do().
		Into(result)
	return
}

// Delete takes name of the qperf and deletes it. Returns an error if one occurs.
func (c *qperves) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("qperves").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *qperves) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("qperves").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched qperf.
func (c *qperves) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *perfv1alpha1.Qperf, err error) {
	result = &perfv1alpha1.Qperf{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("qperves").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
