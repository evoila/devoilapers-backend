package common_test

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/service"
	"errors"
	"time"
)

func WaitForServiceComeUp(provider provider.IServiceProvider, user common.IKubernetesAuthInformation, serviceId string) (*service.IService, error) {
	for i := 0; i < 100; i++ {
		time.Sleep(5 * time.Second)

		// Try get service with invalid user data
		servicePtr, err := provider.GetService(user, serviceId)
		if err != nil {
			return nil, err
		}

		if (*servicePtr).GetStatus() == service.ServiceStatusOk {
			return servicePtr, nil
		}
	}

	return nil, errors.New("Service with id " + serviceId + " does not become ready.")
}

func GetAction(service *service.IService, groupname string, actioncommand string) (*action.IAction, error) {
	actionGroups := (*service).GetActionGroups()

	for _, actionGroup := range actionGroups {
		if actionGroup.GetName() != groupname {
			continue
		}

		actions := actionGroup.GetActions()
		for actionIdx, action := range actions {
			if action.GetUniqueCommand() == actioncommand {
				return &actions[actionIdx], nil
			}
		}
	}

	return nil, errors.New("Action not found")
}



func GetToggleAction(service *service.IService, groupname string, actioncommand string) (*ToggleActionHelper, error) {
	actionPtr, err := GetAction(service, groupname, actioncommand)

	return &ToggleActionHelper{
		Get: func() (bool, error) {
			toggleAction := *actionPtr
			result, err := toggleAction.GetActionExecuteCallback()(&action.ToggleActionPlaceholder{Toggle: "get"})
			return result == true, err
		},
		Set: func() (interface{}, error) {
			toggleAction := *actionPtr
			return toggleAction.GetActionExecuteCallback()(&action.ToggleActionPlaceholder{Toggle: "set"})
		},
		Unset: func() (interface{}, error) {
			toggleAction := *actionPtr
			return toggleAction.GetActionExecuteCallback()(&action.ToggleActionPlaceholder{Toggle: "unset"})
		},

	}, err
}

type ToggleActionHelper struct {
	Get   func() (bool, error)
	Set   func() (interface{}, error)
	Unset func() (interface{}, error)
}