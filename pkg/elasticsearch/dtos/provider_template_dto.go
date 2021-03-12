package dtos

type ProviderYamlTemplateDto struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		Version  string `yaml:"version"`
		NodeSets []struct {
			Name   string      `yaml:"name"`
			Count  interface{} `yaml:"count"`
			Config struct {
				NodeStoreAllowMmap bool `yaml:"node.store.allow_mmap"`
			} `yaml:"config"`
		} `yaml:"nodeSets"`
	} `yaml:"spec"`
}