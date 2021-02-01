package kubernetes

import (
	"bytes"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

// Quelle: https://gist.github.com/pytimer/0ad436972a073bb37b8b6b8b474520fc

type KubernetesWrapper struct {
	Clientset *kubernetes.Clientset
	RestConfig *rest.Config
}

func CreateKubernetesWrapper(host string, token string) (KubernetesWrapper, error) {
	restConfig := rest.Config{
		Host: host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
		APIPath:         "/",
	}

	clientset, err := kubernetes.NewForConfig(&restConfig)
	return KubernetesWrapper{Clientset: clientset, RestConfig: &restConfig}, err
}

func newRestClient(restConfig rest.Config, gv schema.GroupVersion) (rest.Interface, error) {
	//restConfig.ContentConfig = resource.UnstructuredPlusDefaultContentConfig()
	//restConfig.GroupVersion = &gv
	//if len(gv.Group) == 0 {
	//	restConfig.APIPath = "/api"
	//} else {
	//	restConfig.APIPath = "/apis"
	//}

	return rest.RESTClientFor(&restConfig)
}

func getObjFromYAML(yamlDeploy string) v1.Deployment {
	d := v1.Deployment{}
	obj := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(yamlDeploy)), 1000)
	obj.Decode(&d)
	return d
}

func (kubWrapper KubernetesWrapper) Apply(serviceyaml string) error {

	d := v1.Deployment{}
	obj := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(serviceyaml)), 1000)
	obj.Decode(&d)



	groupResources, err := restmapper.GetAPIGroupResources(kubWrapper.Clientset.Discovery())
	if err != nil {
		return err
	}
	rm := restmapper.NewDiscoveryRESTMapper(groupResources)
	_ = rm
	//Get some metadata needed to make the REST request.
	gvk := d.GetObjectKind().GroupVersionKind()
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
	mapping, err := rm.RESTMapping(gk, gvk.Version)
	if err != nil {
		return err
	}

//	// Create a client specifically for creating the object.
	restClient, err := newRestClient(*kubWrapper.RestConfig, mapping.GroupVersionKind.GroupVersion())
	if err != nil {
		return err
	}

	_ = restClient

//	// Use the REST helper to create the object in the "default" namespace.
//	restHelper := resource.NewHelper(restClient, mapping)
	//e, err := restHelper.Create("default", false, d.DeepCopy())

	//_ = e
	//_ = err
//
	return nil
}
	//
	//gr, err := restmapper.GetAPIGroupResources(kubWrapper.Clientset.Discovery())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//mapper := restmapper.NewDiscoveryRESTMapper(gr)
	//mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//var dri dynamic.ResourceInterface
	//if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
	//	if unstructuredObj.GetNamespace() == "" {
	//		unstructuredObj.SetNamespace("default")
	//	}
	//	dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
	//} else {
	//	dri = dd.Resource(mapping.Resource)
	//}
	//
	//if _, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
	//	log.Fatal(err)
	//}
	//
	//return nil
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