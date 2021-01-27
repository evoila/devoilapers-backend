package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api/v1"
)

func GenerateKubeConfig(user string, token string) (*v1.Config, error) {
	var clusterConfig = v1.Config{
		APIVersion: "v1",
		Kind: "Config",
		Clusters: []v1.NamedCluster{
			{
				Name: "minikube",
				Cluster: v1.Cluster{
					Server: "https://192.168.99.114:8443",
					CertificateAuthority: "/home/bene/.minikube/ca.crt",
				},
			},
		},
		Contexts: []v1.NamedContext{
			{
				Name: user,
				Context: v1.Context{
					Cluster: "minikube",
					AuthInfo: user,
				},
			},
		},
		CurrentContext: user,
		AuthInfos: []v1.NamedAuthInfo{
			{
				Name: user,
				AuthInfo: v1.AuthInfo{
					Token: token,
				},
			},
		},
	}
	return &clusterConfig, nil
}

func RestConfig()  {
	restConfig := rest.Config{
		Host: "https://192.168.99.114:8443",
		BearerToken: "",
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	}
	//restClient, err := rest.RESTClientFor(&restConfig)
	clientSet, err := kubernetes.NewForConfig(&restConfig)
	list, err := clientSet.AppsV1().Deployments("default").List(metav1.ListOptions{})
	_ = list
	_ = err
}