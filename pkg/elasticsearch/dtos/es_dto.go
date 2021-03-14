package dtos

type ExposeInformation struct {
	IngressName string `json:"ingress_name" example:"my-ingress`
	HostName    string `json:"host_name" example:"myhost.com`
}

type ScaleInformation struct {
	ReplicasCount int32 `json:"replicas_count" example:"2`
}
