package dummy

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/dummy/dtos"
)

type DummyService struct {
	id          string
	status      int
	yaml        string
	serviceType string
	auth        common.IKubernetesAuthInformation
	state		*bool
}

func (es DummyService) GetYamlTemplate() string {
	return es.yaml
}

func (es DummyService) GetType() string {
	return es.serviceType
}

func (es DummyService) GetName() string {
	return es.id
}

func (es DummyService) GetActions() []action.IActionGroup {
	return []action.IActionGroup{
		action.ActionGroup{
			Name: "Dummy Action Group",
			Actions: []action.IAction{
				action.FormAction{
					Name:          "Dummy Action",
					UniqueCommand: "cmd_dummy_action",
					Placeholder:   &dtos.DummyActionDto{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return es.ExecuteDummyAction(i.(*dtos.DummyActionDto))
					},
				},
				action.ToggleAction{
					Name:          "Dummy Toggle",
					UniqueCommand: "cmd_dummy_toggle_action",
					QueryExecuteCallback: func() (bool, error) {
						return *es.state, nil
					},
					SetExecuteCallback: func() (interface{}, error) {
						t := true
						es.state = &t
						return nil, nil
					},
					UnsetExecuteCallback: func() (interface{}, error) {
						t := false
						es.state = &t
						return nil, nil
					},
				},
			},
		},
	}
}

func (es DummyService) ExecuteDummyAction(dto *dtos.DummyActionDto) (interface{}, error) {
	return dto.Dummy, nil
}

func (es DummyService) GetStatus() int {
	return es.status
}
