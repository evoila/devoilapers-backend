package dtos

type FormResponseDto struct {
	Common FormResponseDtoCommon `yaml:"common"`
}

type FormResponseDtoCommon struct {
	ClusterName string `yaml:"cluster_name"`
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
			} `json:"properties"`
		} `json:"common"`
	} `json:"properties"`
}
