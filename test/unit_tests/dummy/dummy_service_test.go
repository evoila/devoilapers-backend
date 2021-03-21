package dummy

import (
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/dummy"
	"OperatorAutomation/pkg/dummy/dtos"
	"OperatorAutomation/test/unit_tests/common_test"
	"encoding/json"
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
	actionGroups := service0.GetActionGroups()
	assert.Equal(t, 1, len(actionGroups))
	actionGroup0 := actionGroups[0]

	// Group has a name
	assert.NotEqual(t, "", actionGroup0.GetName())

	assert.Equal(t, "1", service0.GetYamlTemplate())
}

func Test_DummyProvider_Service_Actions(t *testing.T) {
	_, service0Ptr := CreateProviderAndService(t)
	service0 := *service0Ptr
	assert.Equal(t, "1", service0.GetYamlTemplate())

	// Get action groups
	actionGroups := service0.GetActionGroups()
	assert.Equal(t, 1, len(actionGroups))
	actionGroup0 := actionGroups[0]

	// Group has a name
	assert.NotEqual(t, "", actionGroup0.GetName())
	actions := actionGroup0.GetActions()
	assert.Equal(t, 2, len(actions))
	formAction := actions[0]
	toggleAction := actions[1]

	// Toggle actions
	assert.Equal(t, true, toggleAction.GetIsToggleAction())
	assert.Equal(t, "Dummy Toggle", toggleAction.GetName())
	assert.Equal(t, "cmd_dummy_toggle_action", toggleAction.GetUniqueCommand())
	assert.NotNil(t, toggleAction.GetJsonFormResultPlaceholder())

	// Get
	placeholder := toggleAction.GetJsonFormResultPlaceholder()
	assert.Nil(t, json.Unmarshal([]byte(`{"toggle":"get"}`), placeholder))
	response, err := toggleAction.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.True(t, common_test.CompareJsonObjectWithJsonString(t, `false`, response))

	// Set
	assert.Nil(t, json.Unmarshal([]byte(`{"toggle":"set"}`), placeholder))
	response, err = toggleAction.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.Nil(t, response)

	// Get again
	assert.Nil(t, json.Unmarshal([]byte(`{"toggle":"get"}`), placeholder))
	response, err = toggleAction.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.True(t, common_test.CompareJsonObjectWithJsonString(t, `true`, response))

	// Unset
	assert.Nil(t, json.Unmarshal([]byte(`{"toggle":"unset"}`), placeholder))
	response, err = toggleAction.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.Nil(t, response)

	// Get again
	assert.Nil(t, json.Unmarshal([]byte(`{"toggle":"get"}`), placeholder))
	response, err = toggleAction.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.True(t, common_test.CompareJsonObjectWithJsonString(t, `false`, response))

	// Form action
	assert.NotEqual(t, "", formAction.GetName())
	assert.NotEqual(t, "", formAction.GetUniqueCommand())

	placeholder = formAction.GetJsonFormResultPlaceholder()
	placholderData := placeholder.(*dtos.DummyActionDto)
	placholderData.Dummy = "NewValue"
	response, err = formAction.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.Equal(t, "NewValue", response)
}
