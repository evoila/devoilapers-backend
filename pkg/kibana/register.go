package kibana

import (
	kbCommon "OperatorAutomation/pkg/kibana/common"
	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/kibana/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)


var SchemeGroupVersion = schema.GroupVersion{Group: kbCommon.GroupName, Version: kbCommon.GroupVersion}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func init() {
	AddToScheme(scheme.Scheme)
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&v1.Kibana{},
		&v1.KibanaList{},
	)

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
