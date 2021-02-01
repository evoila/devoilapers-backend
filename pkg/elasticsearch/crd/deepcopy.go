package crd

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (in *Elasticsearch) DeepCopyInto(out *Elasticsearch) {
	out.TypeMeta = v1.TypeMeta{
		Kind:       in.TypeMeta.Kind,
		APIVersion: in.TypeMeta.APIVersion,
	}
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = ElasticsearchStatus{
		Phase:          in.Status.Phase,
		AvailableNodes: in.Status.AvailableNodes,
		Health:         in.Status.Health,
		Version:        in.Status.Version,
	}
}

func (in *ElasticsearchSpec) DeepCopyInto(out *ElasticsearchSpec) {
	out.Image = in.Image
	out.Version = in.Version
	out.ServiceAccountName = in.ServiceAccountName
	//in.Http.DeepCopyInto(&out.Http)
	in.Auth.DeepCopyInto(&out.Auth)
	//if in.SecureSettings != nil {
	//	l := len(in.SecureSettings)
	//	out.SecureSettings = make([]Kibana.SecureSetting, l)
	//	for i, s := range in.SecureSettings {
	//		s.DeepCopyInto(&out.SecureSettings[i])
	//	}
	//}
	if in.NodeSets != nil {
		l := len(in.NodeSets)
		out.NodeSets = make([]NodeSet, l)
		for i, s := range in.NodeSets {
			s.DeepCopyInto(&out.NodeSets[i])
		}
	}
	if in.RemoteClusters != nil {
		l := len(in.RemoteClusters)
		out.RemoteClusters = make([]RemoteCluster, l)
		for i, s := range in.RemoteClusters {
			s.DeepCopyInto(&out.RemoteClusters[i])
		}
	}
	in.PodDisruptionBudget.DeepCopyInto(&out.PodDisruptionBudget)
	in.Transport.DeepCopyInto(&out.Transport)
	in.UpdateStrategy.DeepCopyInto(&out.UpdateStrategy)
}

func (in *Auth) DeepCopyInto(out *Auth) {
	if in.FileRealm != nil {
		l := len(in.FileRealm)
		out.FileRealm = make([]FileRealm, l)
		for i, s := range in.FileRealm {
			out.FileRealm[i] = FileRealm{s.Secret}
		}
	}
	if in.Roles != nil {
		l := len(in.Roles)
		out.Roles = make([]Role, l)
		for i, s := range in.Roles {
			out.Roles[i] = Role{s.Secret}
		}
	}
}

func (in *NodeSet) DeepCopyInto(out *NodeSet) {
	out.Name = in.Name
	out.Count = in.Count
	in.PodTemplate.DeepCopyInto(&out.PodTemplate)
	if in.VolumeClaimTemplates != nil {
		l := len(in.VolumeClaimTemplates)
		out.VolumeClaimTemplates = make([]VolumeClaimTemplate, l)
		for i, s := range in.VolumeClaimTemplates {
			s.DeepCopyInto(&out.VolumeClaimTemplates[i])
		}
	}
}

func (in *VolumeClaimTemplate) DeepCopyInto(out *VolumeClaimTemplate) {
	out.TypeMeta = v1.TypeMeta{
		Kind:       in.TypeMeta.Kind,
		APIVersion: in.TypeMeta.APIVersion,
	}
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.PersistentVolumeClaimSpec.DeepCopyInto(&out.PersistentVolumeClaimSpec)
	in.PersistentVolumeClaimStatus.DeepCopyInto(&out.PersistentVolumeClaimStatus)
}

func (in *RemoteCluster) DeepCopyInto(out *RemoteCluster) {
	out.Name = in.Name
	//out.ElasticsearchRef = Kibana.ElasticsearchRef{
	//	Name:      in.ElasticsearchRef.Name,
	//	Namespace: in.ElasticsearchRef.Namespace,
	//}
}

func (in *PodDisruptionBudget) DeepCopyInto(out *PodDisruptionBudget) {
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

func (in *PodDisruptionBudgetSpec) DeepCopyInto(out *PodDisruptionBudgetSpec) {
	in.LabelSelector.DeepCopyInto(out.LabelSelector)
	out.MaxUnavailable = intstr.IntOrString{
		Type:   in.MaxUnavailable.Type,
		IntVal: in.MaxUnavailable.IntVal,
		StrVal: in.MaxUnavailable.StrVal,
	}
	out.MinUnavailable = intstr.IntOrString{
		Type:   in.MinUnavailable.Type,
		IntVal: in.MinUnavailable.IntVal,
		StrVal: in.MinUnavailable.StrVal,
	}
}

func (in *Transport) DeepCopyInto(out *Transport) {
	in.Service.DeepCopyInto(&out.Service)
}

func (in *UpdateStrategy) DeepCopyInto(out *UpdateStrategy) {
	out.Version = in.Version
	out.ChangeBudget = ChangeBudget{
		MaxSurge:       in.ChangeBudget.MaxSurge,
		MaxUnavailable: in.ChangeBudget.MaxUnavailable,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *Elasticsearch) DeepCopyObject() runtime.Object {
	out := Elasticsearch{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *ElasticsearchList) DeepCopyObject() runtime.Object {
	out := ElasticsearchList{}
	out.TypeMeta = v1.TypeMeta{
		Kind:       in.TypeMeta.Kind,
		APIVersion: in.TypeMeta.APIVersion,
	}
	in.ListMeta.DeepCopyInto(&out.ListMeta)

	if in.Items != nil {
		out.Items = make([]Elasticsearch, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
