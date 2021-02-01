package Kibana

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (in *Kibana) DeepCopyInto(out *Kibana) {
	out.TypeMeta = v1.TypeMeta{
		Kind:       in.TypeMeta.Kind,
		APIVersion: in.TypeMeta.APIVersion,
	}
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = KibanaSpec{
		Version: in.Spec.Version,
		Count:   in.Spec.Count,
		ElasticsearchRef: ElasticsearchRef{
			Name:      in.Spec.ElasticsearchRef.Name,
			Namespace: in.Spec.ElasticsearchRef.Namespace,
		},
		Image:          in.Spec.Image,
		ServiceAccount: in.Spec.ServiceAccount,
	}
	in.Spec.PodTemplate.DeepCopyInto(&out.Spec.PodTemplate)
	if in.Spec.SecureSettings != nil {
		out.Spec.SecureSettings = make([]SecureSetting, len(in.Spec.SecureSettings))
		for i := range in.Spec.SecureSettings {
			in.Spec.SecureSettings[i].DeepCopyInto(&out.Spec.SecureSettings[i])
		}
	}
	in.Spec.Http.DeepCopyInto(&out.Spec.Http)
	out.Status = KibanaStatus{
		AssociationStatus: in.Status.AssociationStatus,
		AvailableNodes:    in.Status.AvailableNodes,
		Health:            in.Status.Health,
		Version:           in.Status.Version,
	}
}

func (in *SecureSetting) DeepCopyInto(out *SecureSetting) {
	out.Secretname = in.Secretname
	if in.Entries != nil {
		l := len(in.Entries)
		out.Entries = make([]Entry, l)
		for i := 0; i < l; i++ {
			out.Entries[i] = Entry{
				Key:  in.Entries[i].Key,
				Path: in.Entries[i].Path,
			}
		}
	}
}

func (in *Http) DeepCopyInto(out *Http) {
	in.Service.DeepCopyInto(&out.Service)
	in.Tls.DeepCopyInto(&out.Tls)
}

func (in *TLS) DeepCopyInto(out *TLS) {
	out.Cert = Certificate{in.Cert.Secret}
	in.SelfSignedCert.DeepCopyInto(&out.SelfSignedCert)
}

func (in *SelfSignedCertificate) DeepCopyInto(out *SelfSignedCertificate) {
	out.Disabled = in.Disabled
	if in.SubjectAltNames != nil {
		l := len(in.SubjectAltNames)
		out.SubjectAltNames = make([]SubjectAltName, l)
		for i := 0; i < l; i++ {
			out.SubjectAltNames[i] = SubjectAltName{
				Dns: in.SubjectAltNames[i].Dns,
				Ip:  in.SubjectAltNames[i].Ip,
			}
		}
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *Kibana) DeepCopyObject() runtime.Object {
	out := Kibana{}
	in.DeepCopyInto(&out)
	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *KibanaList) DeepCopyObject() runtime.Object {
	out := KibanaList{}
	out.TypeMeta = v1.TypeMeta{
		Kind:       in.TypeMeta.Kind,
		APIVersion: in.TypeMeta.APIVersion,
	}
	in.ListMeta.DeepCopyInto(&out.ListMeta)

	if in.Items != nil {
		out.Items = make([]Kibana, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
	return &out
}
