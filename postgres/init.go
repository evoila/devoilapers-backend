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

func GetPostgresService(k8sApi *K8s.K8sApi, postgresApi *Postgresql.PostgresqlApi, namespace, name string) (*PostgresService, error) {
	service := PostgresService{
		Service: internal.Service{
			K8sApi: k8sApi,
		},
		postgresApi: postgresApi,
	}

	cert, key, err := k8sApi.GetCert(namespace, "pgo.tls")
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
	user, password, err := getUserPassword(k8sApi, namespace, "pgouser-admin")
	if err != nil {
		return nil, err
	}
	service.Comm = &Communication.Comm{
		Client:   client,
		Path:     fmt.Sprintf("%s/%s-%s-http", namespace, name, "pg"),
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
	user := string(secret.Data["username"])
	password := string(secret.Data["password"])
	return user, password, nil
}
