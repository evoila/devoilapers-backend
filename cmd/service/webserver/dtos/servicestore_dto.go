package dtos

type ServiceStoreOverviewDto struct {
	ServiceStoreItems []ServiceStoreItemDto `json:"services"`
}

type ServiceStoreItemDto struct {
	Type        string `json:"type" example:"kibana"`
	Description string `json:"description" example:"Dashboard for elasticsearch"`
	ImageSource string `json:"imageSource" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAIAAA/hNdgD"`
}

type ServiceStoreItemYamlDto struct {
	Yaml string `json:"yaml" example:"item: text"`
}


type ServiceStoreItemFormDto struct {
	FormJson string `json:"formJson" example:"{\"MyJsonObj\": \"Value\"}"`
}