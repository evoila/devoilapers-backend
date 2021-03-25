package kubernetes

import (
	"OperatorAutomation/pkg/utils/logger"
	"bytes"
	"context"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/sirupsen/logrus"

	"encoding/json"

	v1 "k8s.io/api/core/v1"
	v1Beta "k8s.io/api/extensions/v1beta1"
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
	logger.RTrace("Get kubernetes clientset from rest config.")

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

func (api *K8sApi) Apply(b []byte) ([]*unstructured.Unstructured, error) {
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(b), 256)

	// Loop to allow multi document yaml
	result := []*unstructured.Unstructured{}
	for {
		logger.RTrace("Begin parsing to apply yaml in cluster")

		var rawObj runtime.RawExtension
		if err := decoder.Decode(&rawObj); err != nil {
			// Multi document yaml has finished
			if err.Error() == "EOF" {
				logger.RTrace("End of yaml reached")
				return result, nil
			}

			logger.RError(err,"Yaml decoder produced an error")
			return result, err
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if obj == nil {
			return result, err
		}
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			logrus.Error(err)
			return  result, err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(api.ClientSet.Discovery())
		if err != nil {
			logger.RError(err,"Could not resolve api group resources.")
			return result, err
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			logger.RError(err,"Could not identify a preferred resource mapping.")
			return result, err
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
		patchResult, err := dri.Patch(context.Background(), unstructuredObj.GetName(),
			types.ApplyPatchType, data, metav1.PatchOptions{
				FieldManager: "field-manager",
				Force:        &force,
			})

		if err != nil {
			logger.RError(err, "Could not apply patch.")
			return result, err
		}

		result = append(result, patchResult)
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

// create a new tls certificate associated with the provided CRD info
// tlsCert must contain ca.crt, tls.crt and tls.key
func (api *K8sApi) CreateTlsSecretWithoutOwner(secretName string, namespace string, tlsCert map[string][]byte) (string, error) {
	_ = api.ClientSet.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})

	secret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: tlsCert,
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

// Delete a secret based on provided name and namespace
func (api *K8sApi) DeleteSecret(namespace, name string) error {
	return api.ClientSet.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}


func (api *K8sApi) CreateIngressWithHttpsBackend(ingressName string, namespace string, hostname string, tlsSecret string, serviceName string, servicePort int, ) (*v1Beta.Ingress, error) {

	new_ingress  := &v1Beta.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressName,
			Namespace: namespace,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "nginx",
				"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS",
			},
		},
		Spec: v1Beta.IngressSpec{
			TLS: []v1Beta.IngressTLS{
				{
					Hosts:      []string{hostname},
					SecretName: tlsSecret,
				},
			},
			Rules: []v1Beta.IngressRule{
				{
					Host: hostname,
					IngressRuleValue: v1Beta.IngressRuleValue{
						HTTP: &v1Beta.HTTPIngressRuleValue{
							Paths: []v1Beta.HTTPIngressPath{
								{
									Path: "/",
									Backend: v1Beta.IngressBackend{
										ServiceName: serviceName,
										ServicePort: intstr.IntOrString{
											Type:   0,
											IntVal: int32(servicePort),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return api.V1beta1Client.Ingresses(namespace).Create(context.TODO(), new_ingress, metav1.CreateOptions{})
}