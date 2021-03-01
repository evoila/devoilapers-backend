package kubernetes

import (
	"context"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"

	"bytes"

	"encoding/json"

	appV1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
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

// pass byte[] in param instead of filePath (because we dont save the file)
func (api *K8sApi) ApplyFile(filePath, opCode string) (*unstructured.Unstructured, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return api.executeYaml(b, opCode)
}


func (api *K8sApi) Apply(yaml []byte) (*unstructured.Unstructured, error) {
	return api.executeYaml(yaml, "apply")
}

func (api *K8sApi) executeYaml(b []byte, opCode string) (*unstructured.Unstructured, error) {
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
			logrus.Error(err)
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

		if opCode == "create" {
			unstructured, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{
				FieldManager: "field-manager",
			})
			if err != nil {
				logrus.Error(err)
			}
			return unstructured, err
		} else if opCode == "apply" {
			data, _ := json.Marshal(obj)
			force := true
			unstructured, err := dri.Patch(context.Background(), unstructuredObj.GetName(),
				types.ApplyPatchType, data, metav1.PatchOptions{
					FieldManager: "field-manager",
					Force:        &force,
				})
			if err != nil {
				logrus.Error(err)
			} else {
				logrus.Info((unstructured.Object["spec"].(map[string]interface{}))["version"])
			}
			return unstructured, err
		} else if opCode == "delete" {
			err = dri.Delete(context.Background(), unstructuredObj.GetName(), metav1.DeleteOptions{})
			if err != nil {
				logrus.Error(err)
			} else {
				logrus.Info("Delete successfully")
			}
			return nil, err
		}
	}
}

func (api *K8sApi) GetPod(namespace, podname string) (*v1.Pod, error) {
	return api.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podname, metav1.GetOptions{})
}

func (api *K8sApi) GetService(namespace, servicename string) (*v1.Service, error) {
	return api.ClientSet.CoreV1().Services(namespace).Get(context.TODO(), servicename, metav1.GetOptions{})
}

func (api *K8sApi) GetServiceAccount(namespace, name string) (*v1.ServiceAccount, error) {
	return api.ClientSet.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (api *K8sApi) GetReplicaset(namespace, name string) (*appV1.ReplicaSet, error) {
	return api.ClientSet.AppsV1().ReplicaSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (api *K8sApi) GetStatefulset(namespace, name string) (*appV1.StatefulSet, error) {
	return api.ClientSet.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (api *K8sApi) GetDeployment(namespace, name string) (*appV1.Deployment, error) {
	return api.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (api *K8sApi) GetSecret(namespace, name string) (*v1.Secret, error) {
	return api.ClientSet.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (api *K8sApi) GetDeploymentScale(namespace, name string) (*autoscalingv1.Scale, error) {
	return api.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), name, metav1.GetOptions{})
}

func (api *K8sApi) GetStatefulSetScale(namespace, name string) (*autoscalingv1.Scale, error) {
	return api.ClientSet.AppsV1().StatefulSets(namespace).GetScale(context.TODO(), name, metav1.GetOptions{})
}

func (api *K8sApi) CreateSecret(namespace, name string, cert, key []byte) (*v1.Secret, error) {
	new_secret := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       KindSecret,
			APIVersion: ApiversionV1,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"tls.crt": cert,
			"tls.key": key,
		},
		Type: "kubernetes.io/tls",
	}
	return api.ClientSet.CoreV1().Secrets(namespace).Create(context.TODO(), &new_secret, metav1.CreateOptions{})
}

func (api *K8sApi) GetClientCertFromSecret(namespace, secretname string) ([]byte, []byte, error) {
	secret, err := api.GetSecret(namespace, secretname)
	if err != nil {
		return nil, nil, err
	}
	return secret.Data["tls.crt"], secret.Data["tls.key"], nil
}

func (api *K8sApi) CreateServiceAccount(namespace, name string) (*v1.ServiceAccount, error) {
	new := v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       KindServiceaccount,
			APIVersion: ApiversionV1,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return api.ClientSet.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), &new, metav1.CreateOptions{})
}

func (api *K8sApi) PollPods(namespace string) (*v1.PodList, error) {
	return api.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (api *K8sApi) PollServices(namespace string) (*v1.ServiceList, error) {
	return api.ClientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (api *K8sApi) PollReplicaSets(namespace string) (*appV1.ReplicaSetList, error) {
	return api.ClientSet.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (api *K8sApi) PollStatefulsets(namespace string) (*appV1.StatefulSetList, error) {
	return api.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (api *K8sApi) PollDeployments(namespace string) (*appV1.DeploymentList, error) {
	return api.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (api *K8sApi) PollServiceAccounts(namespace string) (*v1.ServiceAccountList, error) {
	return api.ClientSet.CoreV1().ServiceAccounts(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (api *K8sApi) DeleteServiceAccount(namespace, name string) error {
	return api.ClientSet.CoreV1().ServiceAccounts(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (api *K8sApi) DeleteService(namespace, name string) error {
	return api.ClientSet.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (api *K8sApi) UpdateScaleDeployment(namespace, name string, num int32) (*autoscalingv1.Scale, error) {
	scale, err := api.GetDeploymentScale(namespace, name)
	if err != nil {
		return nil, err
	}
	scale.Spec.Replicas = num
	return api.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), name, scale, metav1.UpdateOptions{})
}

func (api *K8sApi) UpdateScaleStatefulSet(namespace, name string, num int32) (*autoscalingv1.Scale, error) {
	scale, err := api.GetStatefulSetScale(namespace, name)
	if err != nil {
		return nil, err
	}
	scale.Spec.Replicas = num
	return api.ClientSet.AppsV1().StatefulSets(namespace).UpdateScale(context.TODO(), name, scale, metav1.UpdateOptions{})
}
