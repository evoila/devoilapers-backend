package dummy

type auth struct {

}

func (a auth) GetKubernetesAccessToken() string {
	panic("implement me")
}

func (a auth) GetKubernetesNamespace() string {
	panic("implement me")
}
