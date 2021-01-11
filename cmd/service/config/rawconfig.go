package config

type RawConfig struct {
	LogLevel                string                  `json:"log_level"`
	Port                    int                     `json:"port"`
	WebserverSllCertificate WebserverSllCertificate `json:"webserver_sll_certificate"`
	Kubernetes              Kubernetes              `json:"kubernetes"`
	User                    []User                  `json:"user"`
}

type WebserverSllCertificate struct {
	PrivateKeyFilePath string `json:"private_key_file_path"`
	PublicKeyFilePath  string `json:"public_key_file_path"`
}

type Kubernetes struct {
	Server               string `json:"server"`
	CertificateAuthority string `json:"certificate-authority"`
	Name                 string `json:"name"`
}

type User struct {
	Name                  string `json:"name"`
	Password              string `json:"password"`
	KubernetesAccessToken string `json:"kubernetes_access_token"`
	KubernetesNamespace   string `json:"kubernetes_namespace"`
	Role                  string `json:"role"`
}

func (user User) GetName() string {
	return user.Name
}
func (user User) GetPassword() string {
	return user.Password
}
func (user User) GetKubernetesAccessToken() string {
	return user.KubernetesAccessToken
}
func (user User) GetKubernetesNamespace() string {
	return user.KubernetesNamespace
}
func (user User) GetRole() string {
	return user.Role
}
