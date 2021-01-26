package kibana

import (
	"crypto/tls"
	"fmt"
	"goclient/Communication"
	"goclient/K8s"
	"goclient/crd/Kibana"
	"goclient/internal"
	"net/http"
)

type KibanaService struct {
	kibanaApi *Kibana.KibanaApi
	internal.Service
}

func GetKibanaService(k8sApi *K8s.K8sApi, kibanaApi *Kibana.KibanaApi, namespace, name string) (*KibanaService, error) {
	service := KibanaService{
		Service: internal.Service{
			K8sApi: k8sApi,
		},
		kibanaApi: kibanaApi,
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

func (service *KibanaService) getCommunication(namespace, name string) (*Communication.Comm, error) {
	kibana, err := service.kibanaApi.Get(namespace, name)
	if err != nil {
		return nil, err
	}
	secretName := kibana.Spec.Http.Tls.Cert.Secret
	var client http.Client
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
	user, password, err := service.getUserPassword(namespace, kibana.Spec.ElasticsearchRef.Name+"-es-elastic-user")
	if err != nil {
		return nil, err
	}
	return &Communication.Comm{
		Client:   &client,
		Path:     fmt.Sprintf("%s/%s-%s-http", namespace, name, "kb"),
		User:     user,
		Password: password,
	}, nil
}

func (service *KibanaService) getUserPassword(namespace, secretName string) (string, string, error) {
	secret, err := service.K8sApi.GetSecret(namespace, secretName)
	if err != nil {
		return "", "", err
	}
	user := "elastic"
	password := string(secret.Data[user])
	return user, password, nil
}

func (service *KibanaService) getAllHandlers() map[string]func([]byte) (string, error) {
	return map[string]func([]byte) (string, error){
		internal.SCALE:        service.handleScale,
		internal.GET_FEATURES: service.handleGetFeatures,
	}
}

func (service *KibanaService) getAllPlaceHolders() map[string]interface{} {
	return map[string]interface{}{
		internal.SCALE: getScalePlaceHolder(),
	}
}
