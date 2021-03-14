package dtos

type FormResponseDto struct {
	Common FormResponseDtoCommon `json:"common"`
}

type FormResponseDtoCommon struct {
	ClusterName           string `json:"cluster_name"`
	ElasticSearchInstance string `json:"elastic_search_instance"`
}

type OneOf struct {
	Description string   `json:"description"`
	Enum        []string `json:"enum"`
}

type FormQueryDto struct {
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
				} `json:"cluster_name"`
				ElasticSearchInstance struct {
					Type   string `json:"type"`
					Title  string `json:"title"`
					Widget struct {
						ID string `json:"id"`
					} `json:"widget"`
					OneOf []OneOf `json:"oneOf"`
				} `json:"elastic_search_instance"`
			} `json:"properties"`
		} `json:"common"`
	} `json:"properties"`
}
