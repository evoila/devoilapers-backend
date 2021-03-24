package yaml_types

type YamlCaSecret struct {
	APIVersion string   `yaml:"apiVersion"`
	Data       CaData  `yaml:"data"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Type       string   `yaml:"type"`
}

type YamlTlsSecret struct {
	APIVersion string   `yaml:"apiVersion"`
	Data       TlsData  `yaml:"data"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Type       string   `yaml:"type"`
}

type CaData struct {
	CaCrtBase64 string `yaml:"ca.crt"`
}

type TlsData struct {
	TLSCrtBase64 string `yaml:"tls.crt"`
	TLSKeyBase64 string `yaml:"tls.key"`
}

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}
