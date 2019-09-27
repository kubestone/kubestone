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

// FiosGetter has a method to return a FioInterface.
// A group's client should implement this interface.
type FiosGetter interface {
	Fios(namespace string) FioInterface
}

// FioInterface has methods to work with Fio resources.
type FioInterface interface {
	Create(*perfv1alpha1.Fio) (*perfv1alpha1.Fio, error)
	Update(*perfv1alpha1.Fio) (*perfv1alpha1.Fio, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*perfv1alpha1.Fio, error)
	List(opts metav1.ListOptions) (*perfv1alpha1.FioList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *perfv1alpha1.Fio, err error)
	FioExpansion
}

// fios implements FioInterface
type fios struct {
	client rest.Interface
	ns     string
}

// newFios returns a Fios
func newFios(c *PerfV1alpha1Client, namespace string) *fios {
	return &fios{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the fio, and returns the corresponding fio object, and an error if there is any.
func (c *fios) Get(name string, options metav1.GetOptions) (result *perfv1alpha1.Fio, err error) {
	result = &perfv1alpha1.Fio{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("fios").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Fios that match those selectors.
func (c *fios) List(opts metav1.ListOptions) (result *perfv1alpha1.FioList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &perfv1alpha1.FioList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("fios").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested fios.
func (c *fios) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("fios").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a fio and creates it.  Returns the server's representation of the fio, and an error, if there is any.
func (c *fios) Create(fio *perfv1alpha1.Fio) (result *perfv1alpha1.Fio, err error) {
	result = &perfv1alpha1.Fio{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("fios").
		Body(fio).
		Do().
		Into(result)
	return
}

// Update takes the representation of a fio and updates it. Returns the server's representation of the fio, and an error, if there is any.
func (c *fios) Update(fio *perfv1alpha1.Fio) (result *perfv1alpha1.Fio, err error) {
	result = &perfv1alpha1.Fio{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("fios").
		Name(fio.Name).
		Body(fio).
		Do().
		Into(result)
	return
}

// Delete takes name of the fio and deletes it. Returns an error if one occurs.
func (c *fios) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("fios").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *fios) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("fios").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched fio.
func (c *fios) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *perfv1alpha1.Fio, err error) {
	result = &perfv1alpha1.Fio{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("fios").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
