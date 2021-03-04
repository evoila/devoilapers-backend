package kubernetes

const (
	KindSecret         = "Secret"
	KindServiceaccount = "ServiceAccount"
	KindIngress        = "Ingress"
	KindRole           = "Role"
	KindRolebinding    = "RoleBinding"

	ApiversionV1              = "v1"
	ApiversionV1beta1         = "v1beta1"
	ApiversionV1authorization = "rbac.authorization.k8s.io/v1"

	ApigroupAuthorization = "rbac.authorization.k8s.io"
)
