package postgres

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	PostgresV1 "github.com/Crunchydata/postgres-operator/pkg/apis/crunchydata.com/v1"
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


// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group-qualified GroupResource.
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// addKnownTypes adds the set of types defined in this package to the supplied scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&PostgresV1.Pgcluster{},
		&PostgresV1.PgclusterList{},
		&PostgresV1.Pgreplica{},
		&PostgresV1.PgreplicaList{},
		&PostgresV1.Pgpolicy{},
		&PostgresV1.PgpolicyList{},
		&PostgresV1.Pgtask{},
		&PostgresV1.PgtaskList{},
	)

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}