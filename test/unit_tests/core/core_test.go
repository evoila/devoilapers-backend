package core

import (
	"OperatorAutomation/pkg/core"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/test/unit_tests/common_test"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func Test_Core_GetProvider(t *testing.T) {
	var providers []*provider.IServiceProvider
	for i := 0; i < 3; i++ {
		closure := i
		var provider provider.IServiceProvider = common_test.TestProvider{
			GetServiceTypeCb: func() string {
				return "Type" + strconv.Itoa(closure)
			},
		}

		providers = append(providers, &provider)
	}

	core1 := core.CreateCore(providers)
	assert.Equal(t, 3, len(core1.Providers))

	provider1, err := core1.GetProviderByName("Type2")
	assert.Nil(t, err)
	assert.NotNil(t, providers)
	assert.Equal(t, providers[2], provider1)
}

func Test_Core_DuplicateProvider_Panic(t *testing.T) {

	var provider1 provider.IServiceProvider = common_test.TestProvider{
		GetServiceTypeCb: func() string {
			return "1"
		},
	}

	var provider2 provider.IServiceProvider = common_test.TestProvider{
		GetServiceTypeCb: func() string {
			return "1"
		},
	}

	var providers = []*provider.IServiceProvider{
		&provider1,
		&provider2,
	}

	// Gets called on destruction. Ensures there was a panic.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Core did not panic")
		}
	}()

	// Should panic
	_ = core.CreateCore(providers)
}

