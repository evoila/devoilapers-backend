package util

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/utils/provider"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BasicService_Attributes(t *testing.T) {
	basicService := provider.BasicService{
		ProviderType: "ProviderType",
		Name: "Name",
		Status: service.ServiceStatusError,
		Yaml: "Yaml",
	}

	assert.Equal(t, "Yaml",  basicService.GetYamlTemplate())
	assert.Equal(t, service.ServiceStatusError,  basicService.GetStatus())
	assert.Equal(t, "Name",  basicService.GetName())
	assert.Equal(t, "ProviderType",  basicService.GetType())
	assert.Equal(t, []action.IActionGroup{},  basicService.GetActionGroups())
}