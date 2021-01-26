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

type ElasticsearchService struct {
	elasticApi *Elasticsearch.ElasticsearchApi
	internal.Service
}

func GetElasticsearchService(k8sApi *K8s.K8sApi, elasticApi *Elasticsearch.ElasticsearchApi, namespace, name string) (*ElasticsearchService, error) {
	service := ElasticsearchService{
		Service: internal.Service{
			Name:      name,
			Namespace: namespace,
			K8sApi:    k8sApi,
		},
		elasticApi: elasticApi,
	}
	if comm, err := service.getCommunication(namespace, name); err != nil {
		return nil, err
	} else {
		service.ActionHandler = service.getAllHandlers()
		service.PlaceHolders = service.getAllPlaceHolders()
		service.Comm = comm
		return &service, nil
	}

}

func (service *ElasticsearchService) getCommunication(namespace, name string) (*Communication.Comm, error) {
	elastic, err := service.elasticApi.Get(namespace, name)
	if err != nil {
		return nil, err
	}
	secretName := elastic.Spec.Http.Tls.Cert.Secret
	client := &http.Client{}
	if len(secretName) > 0 {
		cert, key, err := service.K8sApi.GetCert(namespace, secretName)
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
	user, password, err := service.getUserPassword(namespace, name+"-es-elastic-user")
	if err != nil {
		return nil, err
	}
	return &Communication.Comm{
		Client:   client,
		Path:     fmt.Sprintf("%s/%s-%s-http", namespace, name, "es"),
		User:     user,
		Password: password,
	}, nil
}

func (service *ElasticsearchService) getUserPassword(namespace, secretName string) (string, string, error) {
	secret, err := service.K8sApi.GetSecret(namespace, secretName)
	if err != nil {
		return "", "", err
	}
	user := "elastic"
	password := string(secret.Data[user])
	return user, password, nil
}

func (service *ElasticsearchService) getAllHandlers() map[string]func([]byte) (string, error) {
	return map[string]func([]byte) (string, error){
		internal.SCALE:   service.handleScale,
		internal.BACK_UP: service.handleBackup,

		// my test action
		internal.GET_CONNECTION: service.handleGetConnection,
	}
}

func (service *ElasticsearchService) getAllPlaceHolders() map[string]interface{} {
	return map[string]interface{}{
		internal.BACK_UP: getBackupPlaceHolder(),
		internal.SCALE:   getScalePlaceHolder(),
	}
}
