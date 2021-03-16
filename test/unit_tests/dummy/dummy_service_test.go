package dummy

import (
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/dummy"
	"OperatorAutomation/pkg/dummy/dtos"
	"OperatorAutomation/test/unit_tests/common_test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func CreateProviderAndService(t *testing.T) (*provider.IServiceProvider, *service.IService) {
	auth := common_test.TestUser{}
	var dummyProvider provider.IServiceProvider = dummy.CreateDummyProvider()

	err := dummyProvider.CreateService(auth, "1")
	assert.Nil(t, err)

	// Get a service. Since the name is random we get it by querying all
	services, err := dummyProvider.GetServices(auth)
	service0 := *services[0]

	return &dummyProvider, &service0
}

func Test_DummyProvider_Service_Attributes(t *testing.T) {
	providerPtr, service0Ptr := CreateProviderAndService(t)
	provider := *providerPtr
	service0 := *service0Ptr

	// Compare properties
	assert.NotEqual(t, "", service0.GetName())
	assert.Equal(t, service.ServiceStatusOk, service0.GetStatus())
	assert.Equal(t, provider.GetServiceType(), service0.GetType())

	// Get action groups
	actionGroups := service0.GetActions()
	assert.Equal(t, 1, len(actionGroups))
	actionGroup0 := actionGroups[0]

	// Group has a name
	assert.NotEqual(t, "", actionGroup0.GetName())

	assert.Equal(t, "1", service0.GetYamlTemplate())
}

func Test_DummyProvider_Service_Actions(t *testing.T) {
	_, service0Ptr := CreateProviderAndService(t)
	service0 := *service0Ptr

	// Get action groups
	actionGroups := service0.GetActions()
	assert.Equal(t, 1, len(actionGroups))
	actionGroup0 := actionGroups[0]

	// Group has a name
	assert.NotEqual(t, "", actionGroup0.GetName())
	actions := actionGroup0.GetActions()
	assert.Equal(t, 2, len(actions))
	action0 := actions[0]

	assert.NotEqual(t, "", action0.GetName())
	assert.NotEqual(t, "", action0.GetUniqueCommand())

	placeholder := action0.GetJsonFormResultPlaceholder()
	placholderData := placeholder.(*dtos.DummyActionDto)
	placholderData.Dummy = "NewValue"
	response, err := action0.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.Equal(t, "NewValue", response)

	assert.Equal(t, "1", service0.GetYamlTemplate())
}
