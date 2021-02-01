package kubernetes

const (
	Host     = "https://192.168.99.114:8443"
	CertPath = "/home/tung/.minikube/profiles/minikube/client.crt"
	KeyPath  = "/home/tung/.minikube/profiles/minikube/client.key"
	CAPath   = "/home/tung/.minikube/ca.crt"

	Kind_Secret         = "Secret"
	Kind_ServiceAccount = "ServiceAccount"
	Kind_Ingress        = "Ingress"
	Kind_Role           = "Role"
	Kind_RoleBinding    = "RoleBinding"

	APIVersion_v1              = "v1"
	APIVersion_v1beta1         = "v1beta1"
	APIVersion_v1Authorization = "rbac.authorization.k8s.io/v1"

	APIGroup_Authorization = "rbac.authorization.k8s.io"
)
