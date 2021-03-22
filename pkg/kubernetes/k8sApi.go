package kubernetes

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	rbac "k8s.io/client-go/kubernetes/typed/rbac/v1"
	"k8s.io/client-go/rest"
)

// K8s Api for using k8s functions
// @Attributes:
// ClientSet + Dif: for K8sUtil functions
// V1beta1Client: for K8sIngress
// RbacClient: for K8sAuth
type K8sApi struct {
	ClientSet     *kubernetes.Clientset
	Dif           dynamic.Interface
	V1beta1Client *v1beta1.ExtensionsV1beta1Client
	RbacClient    *rbac.RbacV1Client
}

// generate an K8s api based on provided token
// @params token (string)
// @return *K8sApi, error
func GenerateK8sApiFromToken(host string, caPath string, token string) (*K8sApi, error) {
	config := &rest.Config{
		Host:        host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAFile: caPath,
		},
	}
	clientSet, dif, err := GetClientSet(config)
	if err != nil {
		return nil, err
	} else {
		api := &K8sApi{
			ClientSet: clientSet,
			Dif:       dif,
		}
		if rbacClient, err := rbac.NewForConfig(config); err != nil {
			return nil, err
		} else {
			api.RbacClient = rbacClient
			if v1beta1Client, err := v1beta1.NewForConfig(config); err != nil {
				return nil, err
			} else {
				api.V1beta1Client = v1beta1Client
				return api, nil
			}
		}
	}
}
