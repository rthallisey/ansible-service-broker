/*
Copyright (c) 2018 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package v1

import (
	scheme "github.com/automationbroker/broker-client-go/client/clientset/versioned/scheme"
	v1 "github.com/automationbroker/broker-client-go/pkg/apis/automationbroker.io/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// JobStatesGetter has a method to return a JobStateInterface.
// A group's client should implement this interface.
type JobStatesGetter interface {
	JobStates(namespace string) JobStateInterface
}

// JobStateInterface has methods to work with JobState resources.
type JobStateInterface interface {
	Create(*v1.JobState) (*v1.JobState, error)
	Update(*v1.JobState) (*v1.JobState, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.JobState, error)
	List(opts meta_v1.ListOptions) (*v1.JobStateList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.JobState, err error)
	JobStateExpansion
}

// jobStates implements JobStateInterface
type jobStates struct {
	client rest.Interface
	ns     string
}

// newJobStates returns a JobStates
func newJobStates(c *AutomationbrokerV1Client, namespace string) *jobStates {
	return &jobStates{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the jobState, and returns the corresponding jobState object, and an error if there is any.
func (c *jobStates) Get(name string, options meta_v1.GetOptions) (result *v1.JobState, err error) {
	result = &v1.JobState{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("jobstates").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of JobStates that match those selectors.
func (c *jobStates) List(opts meta_v1.ListOptions) (result *v1.JobStateList, err error) {
	result = &v1.JobStateList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("jobstates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested jobStates.
func (c *jobStates) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("jobstates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a jobState and creates it.  Returns the server's representation of the jobState, and an error, if there is any.
func (c *jobStates) Create(jobState *v1.JobState) (result *v1.JobState, err error) {
	result = &v1.JobState{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("jobstates").
		Body(jobState).
		Do().
		Into(result)
	return
}

// Update takes the representation of a jobState and updates it. Returns the server's representation of the jobState, and an error, if there is any.
func (c *jobStates) Update(jobState *v1.JobState) (result *v1.JobState, err error) {
	result = &v1.JobState{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("jobstates").
		Name(jobState.Name).
		Body(jobState).
		Do().
		Into(result)
	return
}

// Delete takes name of the jobState and deletes it. Returns an error if one occurs.
func (c *jobStates) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("jobstates").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *jobStates) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("jobstates").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched jobState.
func (c *jobStates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.JobState, err error) {
	result = &v1.JobState{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("jobstates").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
