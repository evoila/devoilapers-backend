package elasticsearch

import (
	"encoding/json"
	"fmt"
	form "goclient/RequestForms/Elasticsearch"
)

func (service *ElasticsearchService) handleScale(body []byte) (string, error) {
	var request form.ScaleRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		return "", err
	}

	if _, err = service.K8sApi.UpdateScaleStatefulSet(service.Namespace, service.Name, request.Replicas); err != nil {
		return "", err
	} else {
		return fmt.Sprintf(`{"msg": "scale %s in namespace %s successfully"}`, service.Namespace, service.Name), nil
	}
}

func (service *ElasticsearchService) handleBackup(body []byte) (string, error) {
	var request form.CreateSnapshotRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		return "", err
	}
	return service.Comm.CreateSnapshot(request)
}

// my test handling function
func (service *ElasticsearchService) handleGetConnection(body []byte) (string, error) {
	fmt.Println(string(body))
	return service.Comm.GetConnection()
}
