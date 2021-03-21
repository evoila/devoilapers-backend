package dtos

//NGX schema form result
type ServiceCreationFormResponseDto struct {
	Common struct {
		ClusterName           string `json:"clusterName"`
		ElasticSearchInstance string `json:"elasticsearchInstance"`
	} `json:"common"`
}

// NGX Schema form
type ServiceCreationFormDto struct {
	Properties struct {
		Common struct {
			Type        string `json:"type"`
			Description string `json:"description"`
			Properties  struct {
				ClusterName struct {
					Type    string `json:"type"`
					Title   string `json:"title"`
					Default string `json:"default"`
					Widget  struct {
						ID string `json:"id"`
					} `json:"widget"`
				} `json:"clusterName"`
				ElasticSearchInstance struct {
					Type   string `json:"type"`
					Title  string `json:"title"`
					Widget struct {
						ID string `json:"id"`
					} `json:"widget"`
					OneOf []OneOfElasticSearchInstance `json:"oneOf"`
				} `json:"elasticsearchInstance"`
			} `json:"properties"`
		} `json:"common"`
	} `json:"properties"`
}

type OneOfElasticSearchInstance struct {
	Description string   `json:"description"`
	Enum        []string `json:"enum"`
}
