package common_test

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/service"
)

type TestService struct {
	GetTypeCb        func() string
	GetNameCb        func() string
	GetTemplateCb    func() service.IServiceTemplate
	GetActionGroupCb func() []action.IActionGroup
	GetStatusCb      func() int
}

func (es TestService) GetType() string {
	return es.GetTypeCb()
}

func (es TestService) GetName() string {
	return es.GetNameCb()
}

func (es TestService) GetActions() []action.IActionGroup {
	return es.GetActionGroupCb()
}

func (es TestService) GetTemplate() service.IServiceTemplate {
	return es.GetTemplateCb()
}

func (es TestService) GetStatus() int {
	return es.GetStatusCb()
}