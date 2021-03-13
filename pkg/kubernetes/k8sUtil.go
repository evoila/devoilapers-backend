package kubernetes

import (
	"bytes"
	"context"

	"github.com/sirupsen/logrus"

	"encoding/json"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

func GetClientSet(Config *rest.Config) (*kubernetes.Clientset, dynamic.Interface, error) {
	if clientSet, err := kubernetes.NewForConfig(Config); err != nil {
		return nil, nil, err
	} else {
		if mDynamic, err := dynamic.NewForConfig(Config); err != nil {
			return nil, nil, err
		} else {
			return clientSet, mDynamic, nil
		}
	}
}

func (api *K8sApi) Apply(b []byte) (*unstructured.Unstructured, error) {
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(b), 100)
	for {
		var rawObj runtime.RawExtension
		if err := decoder.Decode(&rawObj); err != nil {
			return nil, err
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if obj == nil {
			return nil, err
		}
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(api.ClientSet.Discovery())
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			logrus.Error("k8sUtil 1", err)
			return nil, err
		}

		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace("default")
			}
			dri = api.Dif.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = api.Dif.Resource(mapping.Resource)
		}

		data, _ := json.Marshal(obj)
		force := true
		unstructured, err := dri.Patch(context.Background(), unstructuredObj.GetName(),
			types.ApplyPatchType, data, metav1.PatchOptions{
				FieldManager: "field-manager",
				Force:        &force,
			})
		if err != nil {
			logrus.Error("k8sUtil 2", err)
		} else {
			logrus.Info((unstructured.Object["spec"].(map[string]interface{}))["version"])
		}
		return unstructured, err

	}
}

// create a new tls certificate associated with the provided CRD info
// tlsCert must contain ca.crt, tls.crt and tls.key
func (api *K8sApi) CreateTlsSecret(namespace, ownerName, kind, apiVersion, uid string, tlsCert map[string][]byte) (string, error) {
	secretName := ownerName + "-tls-cert"
	_ = api.ClientSet.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	True := true
	secret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         apiVersion,
					Name:               ownerName,
					Kind:               kind,
					Controller:         &True,
					BlockOwnerDeletion: &True,
					UID:                types.UID(uid),
				},
			},
		},
		Data: map[string][]byte{
			"ca.crt":  tlsCert["ca.crt"],
			"tls.crt": tlsCert["tls.crt"],
			"tls.key": tlsCert["tls.key"],
		},
		Type: "Opaque",
	}
	if _, err := api.ClientSet.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{}); err != nil {
		return "", err
	}
	return secretName, nil
}

// Get a secret based on provided name and namespace
func (api *K8sApi) GetSecret(namespace, name string) (*v1.Secret, error) {
	return api.ClientSet.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
