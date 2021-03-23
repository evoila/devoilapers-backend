package util

import (
	provider2 "OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/utils/provider"
	"OperatorAutomation/test/unit_tests/common_test"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func creat_basic_provider(t *testing.T) provider.BasicProvider {
	file, err := ioutil.TempFile(os.TempDir(), "opa.*.txt")

	assert.Nil(t, err)
	defer os.Remove(file.Name())

	_, err = file.Write([]byte("SomeText"));
	assert.Nil(t, err)

	provider := provider.CreateCommonProvider(
		"Host",
		file.Name(),
		file.Name(),
		file.Name(),
		"ProviderType",
		"Description",
		"Image",
	)

	return provider
}

func Test_AbstractBasicProvider_NoThrow_Attributes(t *testing.T) {

	basicProvider := creat_basic_provider(t)
	basicProvider.OnCoreInitialized([]*provider2.IServiceProvider{})

	assert.Equal(t, "Image", basicProvider.GetServiceImage())
	assert.Equal(t, "Description", basicProvider.GetServiceDescription())
	assert.Equal(t, "ProviderType", basicProvider.GetServiceType())
}

func Test_AbstractBasicProvider_Throw_GetService(t *testing.T) {
	basicProvider := creat_basic_provider(t)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Abstract method did not throw")
		}
	}()

	auth := common_test.TestUser{}
	_, _ = basicProvider.GetService(auth, "")
}

func Test_AbstractBasicProvider_Throw_GetServices(t *testing.T) {
	basicProvider := creat_basic_provider(t)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Abstract method did not throw")
		}
	}()

	auth := common_test.TestUser{}
	_, _ = basicProvider.GetServices(auth)
}

func Test_AbstractBasicProvider_Throw_CreateService(t *testing.T) {
	basicProvider := creat_basic_provider(t)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Abstract method did not throw")
		}
	}()

	auth := common_test.TestUser{}
	_ = basicProvider.CreateService(auth, "")
}

func Test_AbstractBasicProvider_Throw_DeleteService(t *testing.T) {
	basicProvider := creat_basic_provider(t)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Abstract method did not throw")
		}
	}()

	auth := common_test.TestUser{}
	_ = basicProvider.DeleteService(auth, "")
}

func Test_AbstractBasicProvider_Throw_YamlTemplate(t *testing.T) {
	basicProvider := creat_basic_provider(t)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Abstract method did not throw")
		}
	}()

	auth := common_test.TestUser{}
	_,_ = basicProvider.GetYamlTemplate(auth, []byte(""))
}

func Test_AbstractBasicProvider_Throw_FormTemplate(t *testing.T) {
	basicProvider := creat_basic_provider(t)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Abstract method did not throw")
		}
	}()

	auth := common_test.TestUser{}
	_,_ = basicProvider.GetJsonForm(auth)
}