package service

type Service interface {
	GetServiceTemplate() (ServiceTemplate)
	ExecuteCustomAction(actionName string, commandDataJson string) (string)
}
