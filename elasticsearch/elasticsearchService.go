package elasticsearch

import (
	"encoding/json"
	"fmt"
	"goclient/Communication"
	form "goclient/RequestForms/Elasticsearch"
	"goclient/crd/Elasticsearch"
	"goclient/internal"
)

type ElasticsearchService struct {
	Comm       *Communication.Comm
	elasticApi *Elasticsearch.ElasticsearchApi
	internal.Service
}

func (service *ElasticsearchService) handleScale(body []byte) (string, error) {
	var request form.ScaleRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		return "", err
	}
	namespace, name := request.Namespace, request.Name
	_, err = service.K8sApi.UpdateScaleStatefulSet(namespace, name, request.Replicas)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`{"msg": "scale %s in namespace %s successfully"}`, namespace, name), nil
}

func (service *ElasticsearchService) handleBackup(body []byte) (string, error) {
	var request form.CreateSnapshotRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		return "", err
	}
	return service.Comm.CreateSnapshot(request)
}

func (service *ElasticsearchService) getAllHandlers() map[string]func([]byte) (string, error) {
	return map[string]func([]byte) (string, error){
		internal.SCALE:   service.handleScale,
		internal.BACK_UP: service.handleBackup,

		// my test action
		internal.GET_CONNECTION: service.handleGetConnection,
	}
}

// my test handling function
func (service *ElasticsearchService) handleGetConnection(body []byte) (string, error) {
	fmt.Println(string(body))
	return service.Comm.GetConnection()
}
