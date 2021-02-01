package kubernetes

import (
	"context"
	"fmt"
	v1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"k8s.io/client-go/rest"

)

func GetV1Beta1Client(config *rest.Config) (*v1beta1.ExtensionsV1beta1Client, error) {
	return v1beta1.NewForConfig(config)
}

// get ingress with provided name in given namespace
func (api *K8sApi) GetIngress(namespace, name string) (*v1.Ingress, error) {
	return api.V1beta1Client.Ingresses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (api *K8sApi) CreateIngress(namespace, ingressName, serviceName, hostname string, servicePort int32) (*v1.Ingress, error) {
	new_ingress := v1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       Kind_Ingress,
			APIVersion: APIVersion_v1beta1,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressName,
			Namespace: namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS",
				"nginx.ingress.kubernetes.io/secure-backends":  "true",
				"ingress.kubernetes.io/ssl-passthrough":        "true",
				"kubernetes.io/ingress.class":                  "nginx",
				"nginx.ingress.kubernetes.io/rewrite-target":   "/$1",
			},
		},
		Spec: v1.IngressSpec{
			TLS: []v1.IngressTLS{
				{
					Hosts:      []string{hostname},
					SecretName: namespace + "-ingress-sercret",
				},
			},
			Rules: []v1.IngressRule{
				{
					Host: hostname,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path: "/" + namespace + "/" + serviceName + "/?(.*)",
									Backend: v1.IngressBackend{
										ServiceName: serviceName,
										ServicePort: intstr.IntOrString{
											Type:   0,
											IntVal: servicePort,
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
	return api.V1beta1Client.Ingresses(namespace).Create(context.TODO(), &new_ingress, metav1.CreateOptions{})
}

// add service to ingress, so that requests, which come to service, should be loaded via ingress
func (api *K8sApi) AddServiceToIngress(namespace, ingressName, serviceName, hostname string, servicePort int32) error {
	ingress, err := api.GetIngress(namespace, ingressName)
	if err != nil {
		_, err = api.CreateIngress(namespace, ingressName, serviceName, hostname, servicePort)
		return err
	}
	existing := api.ExistingServiceInIngress(ingress, serviceName)
	if !existing {
		new_path := v1.HTTPIngressPath{
			Path: "/" + namespace + "/" + serviceName + "/?(.*)",
			Backend: v1.IngressBackend{
				ServiceName: serviceName,
				ServicePort: intstr.IntOrString{
					Type:   0,
					IntVal: servicePort,
				},
			},
		}
		ingress.Spec.Rules[0].HTTP.Paths = append(ingress.Spec.Rules[0].HTTP.Paths, new_path)
		_, err = api.V1beta1Client.Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	}
	return err
}

// delete service from ingress when service is deleted
func (api *K8sApi) DeleteServiceFromIngress(namespace, ingressName, serviceName string) error {
	ingress, err := api.GetIngress(namespace, ingressName)
	if err != nil {
		return err
	}
	l := len(ingress.Spec.Rules)
	if l > 1 {
		fmt.Println(ingress.Spec.Rules[0].HTTP.Paths)
		for i, path := range ingress.Spec.Rules[0].HTTP.Paths {
			if path.Backend.ServiceName == serviceName {
				ingress.Spec.Rules[0].HTTP.Paths[i] = ingress.Spec.Rules[0].HTTP.Paths[l-1]
				ingress.Spec.Rules[0].HTTP.Paths = ingress.Spec.Rules[0].HTTP.Paths[:l]
				break
			}
		}
		fmt.Println(ingress.Spec.Rules[0].HTTP.Paths)
		_, err = api.V1beta1Client.Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	} else if l == 1 {
		err = api.V1beta1Client.Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	}
	return err
}

func (api *K8sApi) ExistingServiceInIngress(ingress *v1.Ingress, serviceName string) bool {
	if len(ingress.Spec.Rules) > 0 {
		rule := ingress.Spec.Rules[0]
		for _, path := range rule.HTTP.Paths {
			if path.Backend.ServiceName == serviceName {
				return true
			}
		}
	}
	return false
}
