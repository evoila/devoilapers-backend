package postgres

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	pgV1 "github.com/Crunchydata/postgres-operator/pkg/apis/crunchydata.com/v1"
)

const GroupName = "crunchydata.com"
const GroupVersion = "v1"
const ResourceName = "pgclusters"

var (
	// SchemeBuilder ...
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme ...
	AddToScheme = SchemeBuilder.AddToScheme
)

// SchemeGroupVersion is the group version used to register these objects.
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

func init() {
	AddToScheme(scheme.Scheme)
}

// addKnownTypes adds the set of types defined in this package to the supplied scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&pgV1.Pgcluster{},
		&pgV1.PgclusterList{},
		&pgV1.Pgreplica{},
		&pgV1.PgreplicaList{},
		&pgV1.Pgpolicy{},
		&pgV1.PgpolicyList{},
		&pgV1.Pgtask{},
		&pgV1.PgtaskList{},
	)

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
