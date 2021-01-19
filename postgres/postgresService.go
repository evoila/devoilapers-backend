package postgres

import (
	"encoding/json"
	"goclient/Communication"
	form "goclient/RequestForms/Postgresql"
	"goclient/crd/Postgresql"
	"goclient/internal"
)

type PostgresService struct {
	Comm        *Communication.Comm
	postgresApi *Postgresql.PostgresqlApi
	internal.Service
}

func stringifyResponse(reponse interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(reponse)
	if err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

func (service *PostgresService) handleScale(body []byte) (string, error) {
	var request form.ScaleRequest
	if err := json.Unmarshal(body, &request); err != nil {
		return "", err
	}
	return stringifyResponse(service.Comm.ScaleCluster(request.Name, &request.ClusterScaleRequest))
}

func (service *PostgresService) handleBackup(body []byte) (string, error) {
	var request form.CreateBackupRequest
	if err := json.Unmarshal(body, &request); err != nil {
		return "", nil
	}
	return stringifyResponse(service.Comm.CreateBackrest(&request.CreateBackrestBackupRequest))
}

func (service *PostgresService) getAllHandlers() map[string]func([]byte) (string, error) {
	return map[string]func([]byte) (string, error){
		internal.SCALE:   service.handleScale,
		internal.BACK_UP: service.handleBackup,
	}
}
