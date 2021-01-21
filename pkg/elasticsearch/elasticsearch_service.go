package elasticsearch

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch/dtos"
)

type ElasticSearchService struct {
	auth common.IKubernetesAuthInformation
}

func (es ElasticSearchService) GetActions() []action.IActionGroup {

	return []action.IActionGroup{

		// Part to add a Action
		action.ActionGroup{
			Name: "Backup & Restore",
			Actions: []action.IAction{
				action.Action{
					Name:        "Backup",
					Placeholder: &dtos.BackupActionDto{AwsS3Path: "<PATH>"},
					ActionExecuteCallback: func(i interface{}) (string, error) {
						return es.ExecuteBackup(i.(*dtos.BackupActionDto))
					},
				},
			},
		},

	}
}

func (es ElasticSearchService) ExecuteBackup(dto *dtos.BackupActionDto) (string, error) {
	// part to add a function
    // return service.Comm.CreateSnapshot(dto)
	return "", nil
}

func (es ElasticSearchService) GetTemplate() service.IServiceTemplate {
	panic("implement me")
}

func (es ElasticSearchService) GetStatus() int {
	panic("implement me")
}
