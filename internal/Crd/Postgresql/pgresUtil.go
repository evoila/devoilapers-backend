package Postgresql

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"OperatorAutomation/internal"
)

type PostgresqlApi struct {
	Client *rest.RESTClient
}

func GeneratePgApiBasedOnToken(token string) (*PostgresqlApi, error) {
	config := &rest.Config{
		Host:        internal.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CertFile: internal.CAPath,
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
		return &PostgresqlApi{restClient}, nil
	}
}

// methods of PostgresqlApi struct ...

// clusters

// get cluster with provided name in given namespace
func (api *PostgresqlApi) GetCluster(namespace, name string) (*Pgcluster, error) {
	result := Pgcluster{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgclusterResourcePlural).
		Name(name).
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

// list all clusters within given namespace
func (api *PostgresqlApi) ListCluster(namespace string) (*PgclusterList, error) {
	result := PgclusterList{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgclusterResourcePlural).
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

// delete cluster with provided name in given namespace
func (api *PostgresqlApi) DeleteCluster(namespace, name string) error {
	err := api.Client.Delete().
		Namespace(namespace).
		Resource(PgclusterResourcePlural).
		VersionedParams(&metav1.DeleteOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).Error()
	return err
}

// policies

// get policy with provided name in given namespace
func (api *PostgresqlApi) GetPolicy(namespace, name string) (*Pgpolicy, error) {
	result := Pgpolicy{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgpolicyResourcePlural).
		Name(name).
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

//list all policies within given namespace
func (api *PostgresqlApi) ListPolicy(namespace string) (*PgpolicyList, error) {
	result := PgpolicyList{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgpolicyResourcePlural).
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

// delete policy with provided name in given namespace
func (api *PostgresqlApi) DeletePolicy(namespace, name string) error {
	err := api.Client.Delete().
		Namespace(namespace).
		Resource(PgpolicyResourcePlural).
		VersionedParams(&metav1.DeleteOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).Error()
	return err
}

// replicas

// get replica with provided name in given namespace
func (api *PostgresqlApi) GetReplica(namespace, name string) (*Pgreplica, error) {
	result := Pgreplica{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgreplicaResourcePlural).
		Name(name).
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

//list all replicas within given namespace
func (api *PostgresqlApi) ListReplica(namespace string) (*PgreplicaList, error) {
	result := PgreplicaList{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgreplicaResourcePlural).
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

// delete replica with provided name in given namespace
func (api *PostgresqlApi) DeleteReplica(namespace, name string) error {
	err := api.Client.Delete().
		Namespace(namespace).
		Resource(PgreplicaResourcePlural).
		VersionedParams(&metav1.DeleteOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).Error()
	return err
}

// tasks

// get task with provided name in given namespace
func (api *PostgresqlApi) GetTask(namespace, name string) (*Pgtask, error) {
	result := Pgtask{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgtaskResourcePlural).
		Name(name).
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

//list all tasks within given namespace
func (api *PostgresqlApi) ListTask(namespace string) (*PgtaskList, error) {
	result := PgtaskList{}
	e := api.Client.Get().
		Namespace(namespace).
		Resource(PgtaskResourcePlural).
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)
	return &result, e
}

// delete task with provided name in given namespace
func (api *PostgresqlApi) DeleteTask(namespace, name string) error {
	err := api.Client.Delete().
		Namespace(namespace).
		Resource(PgtaskResourcePlural).
		VersionedParams(&metav1.DeleteOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).Error()
	return err
}
