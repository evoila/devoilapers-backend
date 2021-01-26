package internal

import (
	"errors"
	"fmt"
	"goclient/Communication"
	"goclient/K8s"
)

type IService interface {
	ExecuteCustomAction(action, namespace, name string, body []byte) (string, error)
	GetActions() []string
	GetPlaceHolder(action string) interface{}
}

type Service struct {
	Name          string
	Namespace     string
	Comm          *Communication.Comm
	K8sApi        *K8s.K8sApi
	ActionHandler map[string]func([]byte) (string, error)
	PlaceHolders  map[string]interface{}
}

func (service *Service) ExecuteCustomAction(action string, body []byte) (string, error) {
	if handler, ok := service.ActionHandler[action]; ok {
		return handler(body)
	} else {
		return "", errors.Unwrap(fmt.Errorf("No action %s available", action))
	}
}

func (service *Service) GetActions() []string {
	actions := make([]string, len(service.ActionHandler))
	for k := range service.ActionHandler {
		actions = append(actions, k)
	}
	return actions
}

func (service *Service) GetPlaceHolder(action string) interface{} {
	return service.PlaceHolders[action]
}
