package dtos


type ServiceYamlDto struct {
	Yaml string `json:"yaml"`
}

type ServiceInstanceDetailsOverviewDto struct {
	Instances []ServiceInstanceDetailsDto `json:"services"`
}

type ServiceInstanceDetailsDto struct {
	Name string `json:"name" example:"my_kibana_instance_1"`
	Id string `json:"id" example:"936DA01F-9ABD-4D9D-80C7-02AF85C822A8"`
	Type string `json:"type" example:"kibana"`
	Status string `json:"status" example:"ok"`
	Namespace string `json:"namespace" example:"user_namespace_42"`
}

