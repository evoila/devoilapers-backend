package config

type RawConfig struct {
	LogLevel                string                  `json:"log_level"`
	Port                    int                     `json:"port"`
	WebserverSllCertificate WebserverSllCertificate `json:"webserver_sll_certificate"`
	Kubernetes              Kubernetes              `json:"kubernetes"`
	Users                   []User                  `json:"users"`
	ResourcesTemplatesPath  string                  `json:"resources_templates_path"`
}

type WebserverSllCertificate struct {
	PrivateKeyFilePath string `json:"private_key_file_path"`
	PublicKeyFilePath  string `json:"public_key_file_path"`
}

type Kubernetes struct {
	Server               string    `json:"server"`
	CertificateAuthority string    `json:"certificate-authority"`
	Nginx                Nginx     `json:"nginx"`
	Operators            Operators `json:"operators"`
}

type Nginx struct {
	Namespace        string `json:"namespace"`
	TcpConfigMapName string `json:"tcp_config_map_name"`
	DeploymentName   string `json:"deployment_name"`
	ContainerName    string `json:"container_name"`
}

type Operators struct {
	Postgres Postgres `json:"postgres"`
}

type Postgres struct {
	PgoUrl      string `json:"pgo_url"`
	PgoUsername string `json:"pgo_username"`
	PgoPassword string `json:"pgo_password"`
	PgoVersion  string `json:"pgo_version"`
	PgoCaPath   string `json:"pgo_ca_path"`
}

type User struct {
	Name                  string `json:"name"`
	Password              string `json:"password"`
	KubernetesAccessToken string `json:"kubernetes_access_token"`
	KubernetesNamespace   string `json:"kubernetes_namespace"`
}

// Interface IKubernetesAuthInformation
func (user User) GetKubernetesAccessToken() string {
	return user.KubernetesAccessToken
}

// Interface IKubernetesAuthInformation
func (user User) GetKubernetesNamespace() string {
	return user.KubernetesNamespace
}
