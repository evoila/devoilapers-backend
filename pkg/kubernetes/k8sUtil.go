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
	log "github.com/sirupsen/logrus"
)

const K8sUtilLogPrefix = "File: k8sUtil.go: "

func GetClientSet(Config *rest.Config) (*kubernetes.Clientset, dynamic.Interface, error) {
	log.Trace(K8sUtilLogPrefix + "Get kubernetes clientset from rest config.")

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

func (api *K8sApi) Apply(b []byte) error {
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(b), 256)

	// Loop to allow multi document yaml
	for {
		log.Trace(K8sUtilLogPrefix + "Begin parsing to apply yaml in cluster")

		var rawObj runtime.RawExtension
		if err := decoder.Decode(&rawObj); err != nil {
			// Multi document yaml has finished
			if err.Error() == "EOF" {
				log.Trace(K8sUtilLogPrefix + "End of yaml reached")
				return nil
			}

			log.Error(K8sUtilLogPrefix + "Yaml decoder produced an error", err)
			return err
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if obj == nil {
			return err
		}
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			logrus.Error(err)
			return  err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(api.ClientSet.Discovery())
		if err != nil {
			log.Error(K8sUtilLogPrefix + "Could not resolve api group resources.", err)
			return err
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Error(K8sUtilLogPrefix + "Could not identify a preferred resource mapping.", err)
			return err
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
		_, err = dri.Patch(context.Background(), unstructuredObj.GetName(),
			types.ApplyPatchType, data, metav1.PatchOptions{
				FieldManager: "field-manager",
				Force:        &force,
			})

		if err != nil {
			log.Error(K8sUtilLogPrefix + "Could not apply patch.", err)
			return err
		}
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
