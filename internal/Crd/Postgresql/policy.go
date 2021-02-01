package Postgresql

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PgpolicyResourcePlural ...
const PgpolicyResourcePlural = "pgpolicies"

// PgpolicySpec ...
// swagger:ignore
type PgpolicySpec struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	SQL       string `json:"sql"`
	Status    string `json:"status"`
}

// Pgpolicy ...
// swagger:ignore
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Pgpolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   PgpolicySpec   `json:"spec"`
	Status PgpolicyStatus `json:"status,omitempty"`
}

// PgpolicyList ...
// swagger:ignore
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PgpolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Pgpolicy `json:"items"`
}

// PgpolicyStatus ...
// swagger:ignore
type PgpolicyStatus struct {
	State   PgpolicyState `json:"state,omitempty"`
	Message string        `json:"message,omitempty"`
}

// PgpolicyState ...
// swagger:ignore
type PgpolicyState string

const (
	// PgpolicyStateCreated ...
	PgpolicyStateCreated PgpolicyState = "pgpolicy Created"
	// PgpolicyStateProcessed ...
	PgpolicyStateProcessed PgpolicyState = "pgpolicy Processed"
)
