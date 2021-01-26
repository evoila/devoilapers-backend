package kibana

import (
	"encoding/json"
	"fmt"
	form "goclient/RequestForms/Kibana"
)

func (service *KibanaService) handleScale(body []byte) (string, error) {
	var request form.ScaleRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		return "", err
	}
	_, err = service.K8sApi.UpdateScaleDeployment(service.Namespace, service.Name, request.Replicas)
	if _, err = service.K8sApi.UpdateScaleDeployment(service.Namespace, service.Name, request.Replicas); err != nil {
		return "", err
	} else {
		return fmt.Sprintf(`{"msg": "scale %s in namespace %s successfully"}`, service.Namespace, service.Name), nil
	}

}

// just for testing purpose
func (service *KibanaService) handleGetFeatures(body []byte) (string, error) {
	return service.Comm.GetFeatures()
}
