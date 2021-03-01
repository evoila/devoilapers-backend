package crd

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
)

// Elastic search api for accessing Elasticsearch custom resource definition
type ElasticsearchApi struct {
	Client *rest.RESTClient
}

// generate an elastic search api based on provided token
func CreateElasticsearchApi(host string, caPath string, token string) (*ElasticsearchApi, error) {
	config := &rest.Config{
		Host:        host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAFile:   caPath,
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
		return &ElasticsearchApi{restClient}, nil
	}
}

// get the elastic search custom resource with provided name in given namespace
func (api *ElasticsearchApi) Get(namespace, name string) (*v1.Elasticsearch, error) {
	result := v1.Elasticsearch{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource("elasticsearches").
		Name(name).
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

// list all the elastic search custom resource in given namespace
func (api *ElasticsearchApi) List(namespace string) (*v1.ElasticsearchList, error) {
	result := v1.ElasticsearchList{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource("elasticsearches").
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

// delete an elastic search custom resource with provided name in given namespace
func (api *ElasticsearchApi) Delete(namespace, name string) error {
	err := api.Client.Delete().
		Namespace(namespace).
		Resource("elasticsearches").
		VersionedParams(&metav1.DeleteOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).Error()
	return err
}
