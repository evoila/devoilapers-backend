package core

import (
	"OperatorAutomation/pkg/core"
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/test/unit_tests/common_test"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func Test_UserContext_Create_Get_Delete(t *testing.T) {
	createServiceCounter := 0
	getServiceCounter := 0
	deleteServiceCounter := 0
	getServicesCounter := 0

	var providers []*service.IServiceProvider
	for i := 0; i < 3; i++ {
		closure := i
		var provider service.IServiceProvider = common_test.TestProvider{
			GetServiceTypeCb: func() string {
				return "Type" + strconv.Itoa(closure)
			},
			CreateServiceCb: func(auth common.IKubernetesAuthInformation, yaml string) error {
				assert.Equal(t, 2, closure)
				assert.Equal(t, "Token", auth.GetKubernetesAccessToken())
				assert.Equal(t, "Namespace", auth.GetKubernetesNamespace())
				createServiceCounter += 1
				return nil
			},
			GetServiceCb: func(auth common.IKubernetesAuthInformation, id string) (*service.IService, error) {
				assert.Equal(t, 2, closure)
				assert.Equal(t, "Token", auth.GetKubernetesAccessToken())
				assert.Equal(t, "Namespace", auth.GetKubernetesNamespace())
				getServiceCounter += 1
				return nil, nil
			},
			DeleteServiceCb: func(auth common.IKubernetesAuthInformation, id string) error {
				assert.Equal(t, 2, closure)
				assert.Equal(t, "Token", auth.GetKubernetesAccessToken())
				assert.Equal(t, "Namespace", auth.GetKubernetesNamespace())
				deleteServiceCounter += 1
				return nil
			},
			GetServicesCb: func(auth common.IKubernetesAuthInformation) ([]*service.IService, error) {
				getServicesCounter += 1
				return []*service.IService{nil, nil}, nil
			},
		}

		providers = append(providers, &provider)
	}

	core1 := core.CreateCore(providers)
	ctx := core1.CrateUserContext(common_test.TestUser{KubernetesAccessToken: "Token", KubernetesNamespace: "Namespace"})

	// Create the service
	_ = ctx.CreateServices("Type2", "_")
	assert.Equal(t, 1, createServiceCounter)
	// Not existing provider error
	err := ctx.CreateServices("X", "_")
	assert.NotNil(t, err)

	// Create the service
	_, _ = ctx.GetService("Type2", "_")
	assert.Equal(t, 1, getServiceCounter)
	// Not existing provider error
	_, err = ctx.GetService("X", "_")
	assert.NotNil(t, err)

	// Delete the service
	_ = ctx.DeleteService("Type2", "_")
	assert.Equal(t, 1, deleteServiceCounter)
	// Not existing provider error
	err = ctx.DeleteService("X", "_")
	assert.NotNil(t, err)

	// Get all services
	services, err := ctx.GetServices()
	assert.Nil(t, err)
	assert.Equal(t, 3, getServicesCounter)
	// Check if core merges the services-lists together
	assert.Equal(t, 3*2, len(services))
}



