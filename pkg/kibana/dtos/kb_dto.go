package dtos

type ExposeInformation struct {
	IngressName string `json:"ingress_name" example:"my-ingress`
}

type ScaleInformation struct {
	ReplicasCount int32 `json:"replicas_count" example:"2`
}
