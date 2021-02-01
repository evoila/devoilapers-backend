package K8s

import (
	"OperatorAutomation/internal"
	"context"

	"github.com/sirupsen/logrus"

	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbac "k8s.io/client-go/kubernetes/typed/rbac/v1"
	"k8s.io/client-go/rest"
)

func GetRbacClient(config *rest.Config) (*rbac.RbacV1Client, error) {
	return rbac.NewForConfig(config)
}

func (auth *K8sApi) GetRoleBinding(namespace, name string) (*rbacV1.RoleBinding, error) {
	return auth.RbacClient.RoleBindings(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (auth *K8sApi) GetRole(namespace, name string) (*rbacV1.Role, error) {
	return auth.RbacClient.Roles(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (auth *K8sApi) AddUserTo(old *rbacV1.RoleBinding, namespace, name string) error {
	new := rbacV1.RoleBinding{
		TypeMeta: old.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:      old.ObjectMeta.Name,
			Namespace: old.ObjectMeta.Namespace,
			SelfLink:  old.ObjectMeta.SelfLink,
		},
		Subjects: old.Subjects,
		RoleRef:  old.RoleRef,
	}
	new.Subjects = append(
		new.Subjects,
		rbacV1.Subject{
			Kind:      internal.Kind_ServiceAccount,
			APIGroup:  internal.APIGroup_Authorization,
			Name:      name,
			Namespace: namespace,
		},
	)
	_, err := auth.RbacClient.RoleBindings(namespace).Update(context.TODO(), &new, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func (auth *K8sApi) DeleteUserFrom(old *rbacV1.RoleBinding, name, namespace, kind string) error {
	newSubjects := old.Subjects
	for i, sub := range newSubjects {
		if sub.Kind == kind && sub.Name == name && sub.Namespace == namespace {
			n := len(newSubjects)
			newSubjects[i] = newSubjects[n-1]
			newSubjects = newSubjects[:n]
			break
		}
	}
	new := rbacV1.RoleBinding{
		TypeMeta: old.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:      old.ObjectMeta.Name,
			Namespace: old.ObjectMeta.Namespace,
			SelfLink:  old.ObjectMeta.SelfLink,
		},
		Subjects: newSubjects,
		RoleRef:  old.RoleRef,
	}
	_, err := auth.RbacClient.RoleBindings(namespace).Update(context.TODO(), &new, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func (auth *K8sApi) CreateRoleBinding(namespace, name, rolename string, account_namespace, account_name []string) (*rbacV1.RoleBinding, error) {
	subjects := []rbacV1.Subject{}
	for i := 0; i < len(account_name); i++ {
		subjects = append(
			subjects,
			rbacV1.Subject{
				Kind:      internal.Kind_ServiceAccount,
				Name:      account_name[i],
				Namespace: account_namespace[i],
			},
		)
	}
	new := rbacV1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       internal.Kind_RoleBinding,
			APIVersion: internal.APIGroup_Authorization,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		RoleRef: rbacV1.RoleRef{
			APIGroup: internal.APIGroup_Authorization,
			Kind:     internal.Kind_Role,
			Name:     rolename,
		},
		Subjects: subjects,
	}
	return auth.RbacClient.RoleBindings(namespace).Create(context.TODO(), &new, metav1.CreateOptions{})
}

func (auth *K8sApi) CreateRole(namespace, name string, apiGroups, resources, verbs [][]string) (*rbacV1.Role, error) {
	rules := []rbacV1.PolicyRule{}
	for i := 0; i < len(verbs); i++ {
		rules = append(
			rules,
			rbacV1.PolicyRule{
				APIGroups: apiGroups[i],
				Resources: resources[i],
				Verbs:     verbs[i],
			},
		)
	}
	new := rbacV1.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       internal.Kind_Role,
			APIVersion: internal.APIGroup_Authorization,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Rules: rules,
	}
	return auth.RbacClient.Roles(namespace).Create(context.TODO(), &new, metav1.CreateOptions{})
}

func (auth *K8sApi) DeleteRole(namespace, name string) error {
	return auth.RbacClient.Roles(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (auth *K8sApi) DeleteRoleBinding(namespace, name string) error {
	return auth.RbacClient.RoleBindings(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (auth *K8sApi) PollRoles(namespace, name string) (*rbacV1.RoleList, error) {
	return auth.RbacClient.Roles(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (auth *K8sApi) PollRoleBindings(namespace, name string) (*rbacV1.RoleBindingList, error) {
	return auth.RbacClient.RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
}
