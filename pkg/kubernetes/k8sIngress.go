package kubernetes

import (
	"context"

	v1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"k8s.io/client-go/rest"
)

const (
	Kind_Ingress       = "Ingress"
	APIVersion_v1beta1 = "v1beta1"
)

// GetV1Beta1Client return a clien that use with ingress
func GetV1Beta1Client(config *rest.Config) (*v1beta1.ExtensionsV1beta1Client, error) {
	return v1beta1.NewForConfig(config)
}

// GetIngress return ingress instance with provided name in given namespace if available
func (api *K8sApi) GetIngress(namespace, name string) (*v1.Ingress, error) {
	return api.V1beta1Client.Ingresses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

//CreateIngress creates new ingress service
func (api *K8sApi) createIngress(namespace, ingressName, serviceName, hostname string, servicePort int32) (string, *v1.Ingress, error) {
	new_ingress := v1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       Kind_Ingress,
			APIVersion: APIVersion_v1beta1,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressName,
			Namespace: namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/backend-protocol":        "HTTPS",
				"nginx.ingress.kubernetes.io/secure-backends":         "true",
				"nginx.ingress.kubernetes.io/ssl-passthrough":         "true",
				"kubernetes.io/ingress.class":                         "nginx",
				"nginx.ingress.kubernetes.io/rewrite-target":          "$2",
				"nginx.ingress.kubernetes.io/affinity":                "cookie",
				"nginx.ingress.kubernetes.io/ssl-redirect":            "false",
				"service.beta.kubernetes.io/do-loadbalancer-hostname": hostname,
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
									Path: "/" + namespace + "/" + serviceName + "(/|$)(.*)",
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
	ingress, err := api.V1beta1Client.Ingresses(namespace).Create(context.TODO(), &new_ingress, metav1.CreateOptions{})
	return (hostname + "/" + namespace + "/" + serviceName), ingress, err
}

// AddServiceToIngress adds service to ingress, so that requests, which come to service, should be loaded via ingress
func (api *K8sApi) AddServiceToIngress(namespace, ingressName, serviceName, hostname string, servicePort int32) (string, error) {
	ingress, err := api.GetIngress(namespace, ingressName)
	var url string
	if err != nil {
		url, _, err = api.createIngress(namespace, ingressName, serviceName, hostname, servicePort)
		return url, err
	}
	existing := api.ExistingServiceInIngress(ingress, serviceName)
	if !existing {
		new_path := v1.HTTPIngressPath{
			Path: "/" + namespace + "/" + serviceName + "(/|$)(.*)",
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
		url = (hostname + "/" + namespace + "/" + serviceName)
	} else {
		url = "already exist in ingress"
	}
	return url, err
}

// DeleteServiceFromIngress deletes service from ingress when service is deleted
func (api *K8sApi) DeleteServiceFromIngress(namespace, ingressName, serviceName string) error {
	ingress, err := api.GetIngress(namespace, ingressName)
	if err != nil {
		return err
	}
	l := len(ingress.Spec.Rules[0].HTTP.Paths)
	if l > 1 {
		for i, path := range ingress.Spec.Rules[0].HTTP.Paths {
			if path.Backend.ServiceName == serviceName {
				ingress.Spec.Rules[0].HTTP.Paths[i] = ingress.Spec.Rules[0].HTTP.Paths[l-1]
				ingress.Spec.Rules[0].HTTP.Paths = ingress.Spec.Rules[0].HTTP.Paths[:l-1]
				break
			}
		}
		_, err = api.V1beta1Client.Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	} else if l == 1 {
		if ingress.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName == serviceName {
			err = api.V1beta1Client.Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
		}
	}
	return err
}

// ExistingServiceInIngress return true if a service is available in an Ingress instance
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
