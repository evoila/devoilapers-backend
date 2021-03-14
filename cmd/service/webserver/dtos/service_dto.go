package dtos

type ServiceYamlDto struct {
	Yaml string `json:"yaml"`
}

type ServiceInstanceDetailsOverviewDto struct {
	Instances []ServiceInstanceDetailsDto `json:"services"`
}

type ServiceInstanceDetailsDto struct {
	Name         string                          `json:"name" example:"my_kibana_instance_1"`
	Type         string                          `json:"type" example:"kibana"`
	Status       string                          `json:"status" example:"ok"`
	ActionGroups []ServiceInstanceActionGroupDto `json:"actionGroups"`
}

type ServiceInstanceActionGroupDto struct {
	GroupName string                     `json:"name" example:"Security"`
	Actions   []ServiceInstanceActionDto `json:"actions"`
}

type ServiceInstanceActionDto struct {
	Name    string `json:"name" example:"Expose service"`
	Command string `json:"command" example:"cmdExpose"`
	Form    string `json:"form" example:"ngx json form valid data"`
}
