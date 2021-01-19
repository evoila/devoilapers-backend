package elasticsearch

import (
	"crypto/tls"
	"fmt"
	"goclient/Communication"
	"goclient/K8s"
	"goclient/crd/Elasticsearch"
	"goclient/internal"
	"net/http"
)

func GetElasticsearchService(k8sApi *K8s.K8sApi, elasticApi *Elasticsearch.ElasticsearchApi, namespace, name string) (*ElasticsearchService, error) {
	service := ElasticsearchService{
		Service: internal.Service{
			K8sApi: k8sApi,
		},
		elasticApi: elasticApi,
	}

	elastic, err := elasticApi.Get(namespace, name)
	if err != nil {
		return nil, err
	}
	secretName := elastic.Spec.Http.Tls.Cert.Secret
	client := &http.Client{}
	if len(secretName) > 0 {
		cert, key, err := k8sApi.GetCert(namespace, secretName)
		if err != nil {
			return nil, err
		}
		pair, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{pair},
				InsecureSkipVerify: true,
			},
		}
	} else {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	user, password, err := getUserPassword(k8sApi, namespace, name+"-es-elastic-user")
	if err != nil {
		return nil, err
	}
	service.Comm = &Communication.Comm{
		Client:   client,
		Path:     fmt.Sprintf("%s/%s-%s-http", namespace, name, "es"),
		User:     user,
		Password: password,
	}
	service.ActionHandler = service.getAllHandlers()
	return &service, nil
}

func getUserPassword(k8sApi *K8s.K8sApi, namespace, secretName string) (string, string, error) {
	secret, err := k8sApi.GetSecret(namespace, secretName)
	if err != nil {
		return "", "", err
	}
	user := "elastic"
	password := string(secret.Data[user])
	return user, password, nil
}
