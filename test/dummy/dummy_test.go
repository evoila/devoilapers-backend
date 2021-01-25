package dummy

import (
	"OperatorAutomation/pkg/dummy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CreateDummyService(t *testing.T) {
	dummyProvider := dummy.CreateDummyProvider()
	assert.Nil(t, dummyProvider.CreateService(auth{},"test: yaml"))
	services := dummyProvider.GetServices(auth{})
	assert.Equal(t, 1,len(services))
	assert.Equal(t, "test: yaml", (*services[0]).GetTemplate().GetYAML())
}

