package service

// Yaml template which mirrors the target service-instance
type IServiceTemplate interface {
	GetYAML() string
	GetImportantSections() []string
}

type ServiceTemplate struct{
	Yaml string
	ImportantSections []string
}

func (st ServiceTemplate) GetYAML() string {
	return st.Yaml
}

func (st ServiceTemplate) GetImportantSections() []string {
	return st.ImportantSections
}

