package service

// Defines an abstraction for an service-instance
type Service interface {
	// Get the template on which the service depends
	GetServiceTemplate() ServiceTemplate
	// Execute a custom command. Like cmd_backup, with the s3 path as a data payload
	ExecuteCustomAction(actionName string, commandDataJson string) string
}
