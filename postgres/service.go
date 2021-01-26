package postgres

import (
	"crypto/tls"
	"fmt"
	"goclient/Communication"
	"goclient/K8s"
	"goclient/crd/Postgresql"
	"goclient/internal"
	"net/http"
)

type PostgresService struct {
	postgresApi *Postgresql.PostgresqlApi
	internal.Service
}

func GetPostgresService(k8sApi *K8s.K8sApi, postgresApi *Postgresql.PostgresqlApi, namespace, name string) (*PostgresService, error) {
	service := PostgresService{
		Service: internal.Service{
			Name:      name,
			Namespace: namespace,
			K8sApi:    k8sApi,
		},
		postgresApi: postgresApi,
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

func (service *PostgresService) getCommunication(namespace, name string) (*Communication.Comm, error) {
	cert, key, err := service.K8sApi.GetCert(namespace, "pgo.tls")
	if err != nil {
		return nil, err
	}
	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{pair},
				InsecureSkipVerify: true,
			},
		},
	}
	user, password, err := service.getUserPassword(namespace, "pgouser-admin")
	if err != nil {
		return nil, err
	}
	return &Communication.Comm{
		Client:   client,
		Path:     fmt.Sprintf("%s/%s-%s-http", namespace, name, "pg"),
		User:     user,
		Password: password,
	}, nil
}

func (service *PostgresService) getUserPassword(namespace, secretName string) (string, string, error) {
	secret, err := service.K8sApi.GetSecret(namespace, secretName)
	if err != nil {
		return "", "", err
	}
	user := string(secret.Data["username"])
	password := string(secret.Data["password"])
	return user, password, nil
}

func (service *PostgresService) getAllHandlers() map[string]func([]byte) (string, error) {
	return map[string]func([]byte) (string, error){
		internal.SCALE:   service.handleScale,
		internal.BACK_UP: service.handleBackup,
	}
}

func (service *PostgresService) getAllPlaceHolders() map[string]interface{} {
	return map[string]interface{}{
		internal.BACK_UP: getBackupPlaceHolder(),
		internal.SCALE:   getScalePlaceHolder(),
	}
}
