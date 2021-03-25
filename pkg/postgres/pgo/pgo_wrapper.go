package pgo

import (
	"OperatorAutomation/pkg/utils/logger"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	msgs "github.com/crunchydata/postgres-operator/pkg/apiservermsgs"
	"io/ioutil"
	"net/http"
	"time"
)

type PgoApi struct {
	credentials msgs.BasicAuthCredentials
	version string
	client *http.Client
}

func CreatePgoApi(apiServerUrl string, apiServerVersion string, caCertificatePath string, username string, password string) *PgoApi {
	caCert, err := ioutil.ReadFile(caCertificatePath)
	if err != nil {
		panic(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	defaultTimeout := 60 * time.Second
	return &PgoApi{
		credentials: msgs.BasicAuthCredentials{
			Username: 	username,
			Password: password,
			APIServerURL: apiServerUrl,
		},
		client: &http.Client{
			Timeout: defaultTimeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
						RootCAs: caCertPool,
				},
			},
		},
		version: apiServerVersion,
	}
}

// Performs a request to the endpoint of the postgres operator
func (api *PgoApi) executeRequest(path string, httpMethod string, request interface{}, response interface{}) error {
	ctx := context.TODO()
	jsonValue, _ := json.Marshal(request)
	url := api.credentials.APIServerURL  + path
	logger.RTrace("%s called...[%s]", path, url)

	req, err := http.NewRequestWithContext(ctx, httpMethod, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(api.credentials.Username, api.credentials.Password)

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Postgres operator api returned invalid invalid status code.")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	// Check operator status message
	operatorStatus := msgs.Status{}
	err = json.Unmarshal(bodyBytes, &operatorStatus)
	if err != nil {
		logger.RError(err, "Unable to decode postgres operator status response")
		return err
	}

	if operatorStatus.Code == msgs.Error {
		return errors.New("Postgres operator api error: " + operatorStatus.Msg)
	}

	// Unmarshal final object
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		logger.RError(err, "Unable to decode postgres operator response")
		return err
	}

	return nil
}

// Create database user
func (api *PgoApi) CreateUser(request *msgs.CreateUserRequest) (msgs.CreateUserResponse, error) {
	var response msgs.CreateUserResponse
	request.ClientVersion = api.version
	err := api.executeRequest("/usercreate", http.MethodPost, request, &response)

	return response, err
}

// Delete database user
func (api *PgoApi) DeleteUser(request *msgs.DeleteUserRequest) (msgs.DeleteUserResponse, error) {
	var response msgs.DeleteUserResponse
	request.ClientVersion = api.version
	err := api.executeRequest("/userdelete", http.MethodPost, request, &response)
	return response, err
}

// Get database users
func (api *PgoApi) GetUsers(request *msgs.ShowUserRequest) (msgs.ShowUserResponse, error) {
	var response msgs.ShowUserResponse
	request.ClientVersion = api.version
	err := api.executeRequest("/usershow", http.MethodPost, request, &response)
	return response, err
}

// Create database backup
func (api *PgoApi) CreateBackup(request *msgs.CreateBackrestBackupRequest) (msgs.CreateBackrestBackupResponse, error) {
	var response msgs.CreateBackrestBackupResponse
	err := api.executeRequest("/backrestbackup", http.MethodPost, request, &response)
	return response, err
}