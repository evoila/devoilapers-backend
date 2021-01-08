package demolib

import (
	"OperatorAutomation/pkg/demolib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_HelloWorldEqualsHelloWorld(t *testing.T) {
	assert.Equal(t, "Hello World", demolib.HelloWorld(), "They should be equal")
}
