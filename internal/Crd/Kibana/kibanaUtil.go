package Kibana

import (
	"context"

	"OperatorAutomation/internal"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// kibana api for accessing kibana custom resource
type KibanaApi struct {
	Client *rest.RESTClient
}

// generate KibanaApi based on provided token
func GenerateKbApiBasedOnToken(token string) (*KibanaApi, error) {
	config := &rest.Config{
		Host:        internal.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CertFile: internal.CertPath,
			KeyFile:  internal.KeyPath,
			CAFile:   internal.CAPath,
		},
	}
	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: GroupName, Version: GroupVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	if restClient, err := rest.UnversionedRESTClientFor(&crdConfig); err != nil {
		return nil, err
	} else {
		return &KibanaApi{restClient}, nil
	}
}

// get kibana custom resource with provided name in given namespace
func (api *KibanaApi) Get(namespace, name string) (*Kibana, error) {
	result := Kibana{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource("kibanas").
		Name(name).
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

// list all kibana custom resource within given namespace
func (api *KibanaApi) List(namespace string) (*KibanaList, error) {
	result := KibanaList{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource("kibanas").
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

// delete kibana custom resource with provided name in given namespace
func (api *KibanaApi) Delete(namespace, name string) error {
	err := api.Client.Delete().
		Namespace(namespace).
		Resource("kibanas").
		VersionedParams(&metav1.DeleteOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).Error()
	return err
}
