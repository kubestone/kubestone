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

// PgbenchesGetter has a method to return a PgbenchInterface.
// A group's client should implement this interface.
type PgbenchesGetter interface {
	Pgbenches(namespace string) PgbenchInterface
}

// PgbenchInterface has methods to work with Pgbench resources.
type PgbenchInterface interface {
	Create(*perfv1alpha1.Pgbench) (*perfv1alpha1.Pgbench, error)
	Update(*perfv1alpha1.Pgbench) (*perfv1alpha1.Pgbench, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*perfv1alpha1.Pgbench, error)
	List(opts metav1.ListOptions) (*perfv1alpha1.PgbenchList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *perfv1alpha1.Pgbench, err error)
	PgbenchExpansion
}

// pgbenches implements PgbenchInterface
type pgbenches struct {
	client rest.Interface
	ns     string
}

// newPgbenches returns a Pgbenches
func newPgbenches(c *PerfV1alpha1Client, namespace string) *pgbenches {
	return &pgbenches{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the pgbench, and returns the corresponding pgbench object, and an error if there is any.
func (c *pgbenches) Get(name string, options metav1.GetOptions) (result *perfv1alpha1.Pgbench, err error) {
	result = &perfv1alpha1.Pgbench{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("pgbenches").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Pgbenches that match those selectors.
func (c *pgbenches) List(opts metav1.ListOptions) (result *perfv1alpha1.PgbenchList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &perfv1alpha1.PgbenchList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("pgbenches").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested pgbenches.
func (c *pgbenches) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("pgbenches").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a pgbench and creates it.  Returns the server's representation of the pgbench, and an error, if there is any.
func (c *pgbenches) Create(pgbench *perfv1alpha1.Pgbench) (result *perfv1alpha1.Pgbench, err error) {
	result = &perfv1alpha1.Pgbench{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("pgbenches").
		Body(pgbench).
		Do().
		Into(result)
	return
}

// Update takes the representation of a pgbench and updates it. Returns the server's representation of the pgbench, and an error, if there is any.
func (c *pgbenches) Update(pgbench *perfv1alpha1.Pgbench) (result *perfv1alpha1.Pgbench, err error) {
	result = &perfv1alpha1.Pgbench{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("pgbenches").
		Name(pgbench.Name).
		Body(pgbench).
		Do().
		Into(result)
	return
}

// Delete takes name of the pgbench and deletes it. Returns an error if one occurs.
func (c *pgbenches) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("pgbenches").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *pgbenches) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("pgbenches").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched pgbench.
func (c *pgbenches) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *perfv1alpha1.Pgbench, err error) {
	result = &perfv1alpha1.Pgbench{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("pgbenches").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
