package kibana

import (
	"encoding/json"
	"errors"
	"fmt"
	"goclient/Communication"
	form "goclient/RequestForms/Kibana"
	"goclient/crd/Kibana"
	"goclient/internal"
)

type KibanaService struct {
	Comm      *Communication.Comm
	kibanaApi *Kibana.KibanaApi
	internal.Service
}

func (service *KibanaService) ExecuteCustomAction(action string, body []byte) (string, error) {
	if handler, ok := service.ActionHandler[action]; ok {
		return handler(body)
	} else {
		return "", errors.Unwrap(fmt.Errorf("No action %s available", action))
	}
}

func (service *KibanaService) handleScale(body []byte) (string, error) {
	var request form.ScaleRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		return "", err
	}
	namespace, name := request.Namespace, request.Name
	_, err = service.K8sApi.UpdateScaleDeployment(namespace, name, request.Replicas)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`{"msg": "scale %s in namespace %s successfully"}`, namespace, name), nil
}

func (service *KibanaService) getAllHandlers() map[string]func([]byte) (string, error) {
	return map[string]func([]byte) (string, error){
		internal.SCALE: service.handleScale,
	}
}
