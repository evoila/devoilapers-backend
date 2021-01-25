package dummy

import (
	"errors"
	"math/rand"
)

// Fake Kubernetes information
type DummyKubernetes struct {
	data map[string]DummyKubernetesData
}
type DummyKubernetesData struct {
	status int
	yaml   string
}

func (dk DummyKubernetes) Create(yaml string) error {
	dk.data[string(rune(rand.Int()))] = DummyKubernetesData{
		status: 1,
		yaml:   yaml,
	}
	return nil
}

func (dk DummyKubernetes) Delete(id string) error {
	delete(dk.data, id)
	return nil
}

func (dk DummyKubernetes) GetServices() map[string]DummyKubernetesData {
	return dk.data
}

func (dk DummyKubernetes) GetService(id string) (DummyKubernetesData, error) {
	data, exists := dk.data[id]
	if exists {
		return data, nil
	}
	return data, errors.New("We dont have this service here!")
}
