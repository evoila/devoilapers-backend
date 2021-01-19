package internal

import (
	"errors"
	"fmt"
	"goclient/K8s"
)

type Service struct {
	K8sApi        *K8s.K8sApi
	ActionHandler map[string]func([]byte) (string, error)
}

func (service *Service) ExecuteCustomAction(action string, body []byte) (string, error) {
	if handler, ok := service.ActionHandler[action]; ok {
		return handler(body)
	} else {
		return "", errors.Unwrap(fmt.Errorf("No action %s available", action))
	}
}

func (service *Service) GetActions(svc string) []string {
	actions := make([]string, len(service.ActionHandler))
	for k := range service.ActionHandler {
		actions = append(actions, k)
	}
	return actions
}

//func (service *Service) GetTemplate(service, action string) interface{}
