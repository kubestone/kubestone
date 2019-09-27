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

// IopingsGetter has a method to return a IopingInterface.
// A group's client should implement this interface.
type IopingsGetter interface {
	Iopings(namespace string) IopingInterface
}

// IopingInterface has methods to work with Ioping resources.
type IopingInterface interface {
	Create(*perfv1alpha1.Ioping) (*perfv1alpha1.Ioping, error)
	Update(*perfv1alpha1.Ioping) (*perfv1alpha1.Ioping, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*perfv1alpha1.Ioping, error)
	List(opts metav1.ListOptions) (*perfv1alpha1.IopingList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *perfv1alpha1.Ioping, err error)
	IopingExpansion
}

// iopings implements IopingInterface
type iopings struct {
	client rest.Interface
	ns     string
}

// newIopings returns a Iopings
func newIopings(c *PerfV1alpha1Client, namespace string) *iopings {
	return &iopings{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the ioping, and returns the corresponding ioping object, and an error if there is any.
func (c *iopings) Get(name string, options metav1.GetOptions) (result *perfv1alpha1.Ioping, err error) {
	result = &perfv1alpha1.Ioping{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("iopings").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Iopings that match those selectors.
func (c *iopings) List(opts metav1.ListOptions) (result *perfv1alpha1.IopingList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &perfv1alpha1.IopingList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("iopings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested iopings.
func (c *iopings) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("iopings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a ioping and creates it.  Returns the server's representation of the ioping, and an error, if there is any.
func (c *iopings) Create(ioping *perfv1alpha1.Ioping) (result *perfv1alpha1.Ioping, err error) {
	result = &perfv1alpha1.Ioping{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("iopings").
		Body(ioping).
		Do().
		Into(result)
	return
}

// Update takes the representation of a ioping and updates it. Returns the server's representation of the ioping, and an error, if there is any.
func (c *iopings) Update(ioping *perfv1alpha1.Ioping) (result *perfv1alpha1.Ioping, err error) {
	result = &perfv1alpha1.Ioping{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("iopings").
		Name(ioping.Name).
		Body(ioping).
		Do().
		Into(result)
	return
}

// Delete takes name of the ioping and deletes it. Returns an error if one occurs.
func (c *iopings) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("iopings").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *iopings) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("iopings").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched ioping.
func (c *iopings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *perfv1alpha1.Ioping, err error) {
	result = &perfv1alpha1.Ioping{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("iopings").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
