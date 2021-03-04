package core

import (
	"OperatorAutomation/pkg/core/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ServiceTemplate_Getter(t *testing.T) {
	var template service.IServiceTemplate = service.ServiceTemplate{
		ImportantSections: []string{"1", "2"},
		Yaml: "Y",
	}

	assert.Equal(t, "1", template.GetImportantSections()[0])
	assert.Equal(t, "2", template.GetImportantSections()[1])
	assert.Equal(t, "Y", template.GetYAML())
}
