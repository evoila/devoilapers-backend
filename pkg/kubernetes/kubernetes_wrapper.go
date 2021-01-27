package kubernetes

import (
	"OperatorAutomation/pkg/core/common"
)

type KubernetesWrapper struct {
}
// Quelle: https://gist.github.com/pytimer/0ad436972a073bb37b8b6b8b474520fc
func CreateKubernetesWrapper(userInformation common.IKubernetesAuthInformation) *KubernetesWrapper {
	return &KubernetesWrapper{}
}

//func Apply(serviceyaml string) error {
//	var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
//
//	obj := &unstructured.Unstructured{}
//	// Decode YAML to unstructured object.
//	if _, _, err := decUnstructured.Decode([]byte(serviceyaml), nil, obj); err != nil {
//		return err
//	}
//	var _, kubeconfig = files.GenerateKubeconfig("","")
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	util
//	config, err := config.GenerateKubeConfig("","")
//
//	restClient := rest.Config{BearerToken: ""}
//
//	clientConfig = client.Config{}
//	clientConfig.Host = "example.com:4901"
//	clientConfig = info.MergeWithConfig()
//	client := client.New(clientConfig)
//	client.Pods(ns).List()
//
//
//	c, err = files.GenerateKubeconfig("","")
//
//	gr, err := restmapper.GetAPIGroupResources(c.Discovery())
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	mapper := restmapper.NewDiscoveryRESTMapper(gr)
//	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	var dri dynamic.ResourceInterface
//	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
//		if unstructuredObj.GetNamespace() == "" {
//			unstructuredObj.SetNamespace("default")
//		}
//		dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
//	} else {
//		dri = dd.Resource(mapping.Resource)
//	}
//
//	if _, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
//		log.Fatal(err)
//	}
//
//	return nil
//}
//
//func (api *K8sApi) Apply(b []byte, opCode string) (*unstructured.Unstructured, error) {
//	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(b), 100)
//	for {
//		var rawObj runtime.RawExtension
//		if err := decoder.Decode(&rawObj); err != nil {
//			return nil, err
//		}
//
//		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
//		if obj == nil {
//			return nil, err
//		}
//		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
//		if err != nil {
//			logrus.Error(err)
//			return nil, err
//		}
//
//		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
//
//		gr, err := restmapper.GetAPIGroupResources(api.ClientSet.Discovery())
//		if err != nil {
//			logrus.Error(err)
//			return nil, err
//		}
//
//		mapper := restmapper.NewDiscoveryRESTMapper(gr)
//		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
//		if err != nil {
//			logrus.Error(err)
//			return nil, err
//		}
//
//		var dri dynamic.ResourceInterface
//		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
//			if unstructuredObj.GetNamespace() == "" {
//				unstructuredObj.SetNamespace("default")
//			}
//			dri = api.Dif.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
//		} else {
//			dri = api.Dif.Resource(mapping.Resource)
//		}
//
//		if opCode == "create" {
//			unstructured, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{
//				FieldManager: "field-manager",
//			})
//			if err != nil {
//				logrus.Error(err)
//			}
//			return unstructured, err
//		} else if opCode == "apply" {
//			data, _ := json.Marshal(obj)
//			force := true
//			unstructured, err := dri.Patch(context.Background(), unstructuredObj.GetName(),
//				types.ApplyPatchType, data, metav1.PatchOptions{
//					FieldManager: "field-manager",
//					Force:        &force,
//				})
//			if err != nil {
//				logrus.Error(err)
//			} else {
//				logrus.Info((unstructured.Object["spec"].(map[string]interface{}))["version"])
//			}
//			return unstructured, err
//		} else if opCode == "delete" {
//			err = dri.Delete(context.Background(), unstructuredObj.GetName(), metav1.DeleteOptions{})
//			if err != nil {
//				logrus.Error(err)
//			} else {
//				logrus.Info("Delete successfully")
//			}
//			return nil, err
//		}
//	}
//}