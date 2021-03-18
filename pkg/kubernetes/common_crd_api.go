package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// Crd api
// Use factory method CreateCommonCrdApi to initialise.
type CommonCrdApi struct {
	Client *rest.RESTClient
}

// Create a kubernetes rest config from parameters
func CreateRestConfig(host string, caPath string, token string, groupName string, groupVersion string) *rest.Config {
	config := &rest.Config{
		Host:        host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAFile: caPath,
		},
	}

	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: groupName, Version: groupVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	return &crdConfig
}

// Create a common crd api to get, list and delete a custom resource in a kubernetes cluster
func CreateCommonCrdApi(host string, caPath string, token string, groupName string, groupVersion string) (*CommonCrdApi, error) {
	crdConfig := CreateRestConfig(host, caPath, token, groupName, groupVersion)

	if restClient, err := rest.UnversionedRESTClientFor(crdConfig); err != nil {
		return nil, err
	} else {
		return &CommonCrdApi{restClient}, nil
	}
}

// Get a single custom resource of given type resource in given namespace with
// given name an pass into given out object
func (api CommonCrdApi) Get(namespace string, name string, resource string, out runtime.Object) error {
	return api.Client.Get().
		Namespace(namespace).
		Resource(resource).
		Name(name).
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(out)
}

// Get a all custom resource of given type resource in given namespace with
// given name an pass into given out (list)-object
func (api CommonCrdApi) List(namespace string, resource string, out runtime.Object) error {
	return api.ListWithOptions(namespace, resource, &metav1.ListOptions{}, out)
}

// Get a all custom resource of given type resource in given namespace with
// given name an pass into given out (list)-object
func (api CommonCrdApi) ListWithOptions(namespace string, resource string, listOptions *metav1.ListOptions, out runtime.Object) error {
	return api.Client.Get().
		Namespace(namespace).
		Resource(resource).
		VersionedParams(listOptions, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(out)
}

// Delete a custom resource of given type resource in given namespace with given name
func (api CommonCrdApi) Delete(namespace string, name string, resource string) error {
	return api.Client.Delete().
		Namespace(namespace).
		Resource(resource).
		VersionedParams(&metav1.DeleteOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).Error()
}

// Update a custom resource of given type resource in given namespace with given name
func (api CommonCrdApi) Update(namespace, name, resource string, obj runtime.Object) error {
	err := api.Client.Put().
		Namespace(namespace).
		Resource(resource).
		Name(name).
		VersionedParams(&metav1.UpdateOptions{}, scheme.ParameterCodec).
		Body(obj).
		Do(context.TODO()).
		Error()
	return err
}
