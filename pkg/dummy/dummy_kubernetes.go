package dummy

import (
	"OperatorAutomation/pkg/core/service"
	"errors"
	"math/rand"
	"strconv"
)

// Fake Kubernetes information
type DummyKubernetes struct {
	data map[string]DummyKubernetesData
}

type DummyKubernetesData struct {
	status      int
	yaml        string
	toggleState *bool
}

func (dk DummyKubernetes) Create(yaml string) error {
	toggleState := false
	dk.data[strconv.Itoa(rand.Int())] = DummyKubernetesData{
		status:      service.ServiceStatusOk,
		yaml:        yaml,
		toggleState: &toggleState,
	}
	return nil
}

func (dk DummyKubernetes) Delete(id string) error {
	if _, found := dk.data[id]; !found {
		return errors.New("Id not found")
	}

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
