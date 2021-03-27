package provider_dtos

type ServiceCreationFormResponseDto struct {
	Common struct {
		ClusterName string `json:"clusterName"`
	} `json:"common"`
}

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
			} `json:"properties"`
		} `json:"common"`
	} `json:"properties"`
}
