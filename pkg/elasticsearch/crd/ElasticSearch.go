package crd

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Elasticsearch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ElasticsearchSpec   `json:"spec,omitempty"`
	Status            ElasticsearchStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ElasticsearchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Elasticsearch `json:"items"`
}

type ElasticsearchSpec struct {
	Auth                Auth `json:"auth,omitempty"`
	// Todo add additional structs
	//Kibana.Http         `json:"http,omitempty"`
	Image               string                 `json:"image,omitempty"`
	Version             string                 `json:"version"`
	NodeSets            []NodeSet              `json:"nodeSets"`
	PodDisruptionBudget PodDisruptionBudget    `json:"podDisruptionBudget,omitempty"`
	RemoteClusters      []RemoteCluster        `json:"remoteClusters,omitempty"`
	//SecureSettings      []Kibana.SecureSetting `json:"secureSettings,omitempty"`
	ServiceAccountName  string                 `json:"serviceAccountName,omitempty"`
	Transport           Transport              `json:"transport,omitempty"`
	UpdateStrategy      UpdateStrategy         `json:"updateStrategy,omitempty"`
}

type Auth struct {
	FileRealm []FileRealm `json:"fileRealm,omitempty"`
	Roles     []Role      `json:"roles,omitempty"`
}

type FileRealm struct {
	Secret string `json:"secretName"`
}

type Role struct {
	Secret string `json:"secretName"`
}

type NodeSet struct {
	// there is also a config for kibana. But its not well defined in crd
	Name                 string `json:"name"`
	Count                int    `json:"count"`
	v1.PodTemplate       `json:"podTemplate,omitempty"`
	VolumeClaimTemplates []VolumeClaimTemplate `json:"volumeClaimTemplates,omitempty"`
}

type VolumeClaimTemplate struct {
	metav1.TypeMeta                `json:",inline"`
	metav1.ObjectMeta              `json:"metadata,omitempty"`
	v1.PersistentVolumeClaimSpec   `json:"spec,omitempty"`
	v1.PersistentVolumeClaimStatus `json:"status,omitempty"`
}

type PodDisruptionBudget struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PodDisruptionBudgetSpec `json:"spec,omitempty"`
}

type PodDisruptionBudgetSpec struct {
	*metav1.LabelSelector `json:"selector,omitempty"`
	MaxUnavailable        intstr.IntOrString `json:"maxUnavailable"`
	MinUnavailable        intstr.IntOrString `json:"minUnavailable"`
}

type RemoteCluster struct {
	Name                    string `json:"name"`
	//Kibana.ElasticsearchRef `json:"elasticsearchRef,omitempty"`
}

type Transport struct {
	v1.Service `json:"service,omitempty"`
}

type UpdateStrategy struct {
	ChangeBudget ChangeBudget `json:"changeBudget,omitempty"`
	Version      string       `json:"version,omitempty"`
}

type ChangeBudget struct {
	MaxSurge       int32 `json:"maxSurge,omitempty"`
	MaxUnavailable int32 `json:"maxUnavailable,omitempty"`
}

type ElasticsearchStatus struct {
	AvailableNodes int32  `json:"availableNodes,omitempty"`
	Health         string `json:"health,omitempty"`
	Phase          string `json:"phase,omitempty"`
	Version        string `json:"version,omitempty"`
}
