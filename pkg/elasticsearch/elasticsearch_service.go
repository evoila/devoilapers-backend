package elasticsearch

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch/dtos"
)

type ElasticSearchService struct {
	serviceType string
	serviceName string
	serviceStatus string
}

func (es ElasticSearchService) GetType() string {
	return es.serviceType
}

func (es ElasticSearchService) GetName() string {
	return es.serviceName
}

func (es ElasticSearchService) GetActions() []action.IActionGroup {

	return []action.IActionGroup{

		// Part to add a Action
		//action.ActionGroup{
		//	Name: "Backup & Restore",
		//	Actions: []action.IAction{
		//		action.Action{
		//			Name:        "Backup",
		//			UniqueCommand: "cmd_elasticsearch_backup",
		//			Placeholder: &dtos.BackupActionDto{},
		//			ActionExecuteCallback: func(i interface{}) (string, error) {
		//				return es.ExecuteBackup(i.(*dtos.BackupActionDto))
		//			},
		//		},
		//	},
		//},

	}
}

func (es ElasticSearchService) ExecuteBackup(dto *dtos.BackupActionDto) (interface{}, error) {
	// part to add a function
    // return service.Comm.CreateSnapshot(dto)
	return "Its OK", nil
}

func (es ElasticSearchService) GetTemplate() service.IServiceTemplate {
	panic("implement me")
}

func (es ElasticSearchService) GetStatus() int {
	if es.serviceStatus == "green" {
		return 2
	} else if es.serviceStatus == "yellow" {
		return 1
	}

	return 0
}
