package provider_dtos

type ProviderYamlTemplateDto struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		Version          string `yaml:"version"`
		Count            int    `yaml:"count"`
		ElasticsearchRef struct {
			Name string `yaml:"name"`
		} `yaml:"elasticsearchRef"`
	} `yaml:"spec"`
}
