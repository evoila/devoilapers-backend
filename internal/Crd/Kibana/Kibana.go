package Kibana

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Kibana struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KibanaSpec   `json:"spec,omitempty"`
	Status            KibanaStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KibanaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kibana `json:"items"`
}

type KibanaSpec struct {
	// there is also a config for kibana. But its not well defined in crd
	Version          string           `json:"version"`
	Count            int32            `json:"count"`
	ElasticsearchRef ElasticsearchRef `json:"elasticsearchRef,omitempty"`
	Http             Http             `json:"http,omitempty"`
	Image            string           `json:"image,omitempty"`
	v1.PodTemplate   `json:"podTemplate,omitempty"`
	SecureSettings   []SecureSetting `json:"secureSettings,omitempty"`
	ServiceAccount   string          `json:"serviceAccountName,omitempty"`
}

type ElasticsearchRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type Http struct {
	v1.Service `json:"service,omitempty"`
	Tls        TLS `json:"tls,omitempty"`
}

type TLS struct {
	Cert           Certificate           `json:"certificate,omitempty"`
	SelfSignedCert SelfSignedCertificate `json:"selfSignedCertificate,omitempty"`
}

type Certificate struct {
	Secret string `json:"secretName,omitempty"`
}

type SelfSignedCertificate struct {
	Disabled        bool             `json:"disabled,omitempty"`
	SubjectAltNames []SubjectAltName `json:"subjectAltNames,omitempty"`
}

type SubjectAltName struct {
	Dns string `json:"dns,omitempty"`
	Ip  string `json:"ip,omitempty"`
}

type SecureSetting struct {
	Entries    []Entry `json:"entries,omitempty"`
	Secretname string  `json:"secretName"`
}

type Entry struct {
	Key  string `json:"key"`
	Path string `json:"path"`
}

type KibanaStatus struct {
	AssociationStatus string `json:"associationStatus,omitempty"`
	AvailableNodes    int32  `json:"availableNodes,omitempty"`
	Health            string `json:"health,omitempty"`
	Version           string `json:"version,omitempty"`
}
