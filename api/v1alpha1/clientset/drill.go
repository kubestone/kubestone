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

// DrillsGetter has a method to return a DrillInterface.
// A group's client should implement this interface.
type DrillsGetter interface {
	Drills(namespace string) DrillInterface
}

// DrillInterface has methods to work with Drill resources.
type DrillInterface interface {
	Create(*perfv1alpha1.Drill) (*perfv1alpha1.Drill, error)
	Update(*perfv1alpha1.Drill) (*perfv1alpha1.Drill, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*perfv1alpha1.Drill, error)
	List(opts metav1.ListOptions) (*perfv1alpha1.DrillList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *perfv1alpha1.Drill, err error)
	DrillExpansion
}

// drills implements DrillInterface
type drills struct {
	client rest.Interface
	ns     string
}

// newDrills returns a Drills
func newDrills(c *PerfV1alpha1Client, namespace string) *drills {
	return &drills{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the drill, and returns the corresponding drill object, and an error if there is any.
func (c *drills) Get(name string, options metav1.GetOptions) (result *perfv1alpha1.Drill, err error) {
	result = &perfv1alpha1.Drill{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("drills").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Drills that match those selectors.
func (c *drills) List(opts metav1.ListOptions) (result *perfv1alpha1.DrillList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &perfv1alpha1.DrillList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("drills").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested drills.
func (c *drills) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("drills").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a drill and creates it.  Returns the server's representation of the drill, and an error, if there is any.
func (c *drills) Create(drill *perfv1alpha1.Drill) (result *perfv1alpha1.Drill, err error) {
	result = &perfv1alpha1.Drill{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("drills").
		Body(drill).
		Do().
		Into(result)
	return
}

// Update takes the representation of a drill and updates it. Returns the server's representation of the drill, and an error, if there is any.
func (c *drills) Update(drill *perfv1alpha1.Drill) (result *perfv1alpha1.Drill, err error) {
	result = &perfv1alpha1.Drill{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("drills").
		Name(drill.Name).
		Body(drill).
		Do().
		Into(result)
	return
}

// Delete takes name of the drill and deletes it. Returns an error if one occurs.
func (c *drills) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("drills").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *drills) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("drills").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched drill.
func (c *drills) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *perfv1alpha1.Drill, err error) {
	result = &perfv1alpha1.Drill{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("drills").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
