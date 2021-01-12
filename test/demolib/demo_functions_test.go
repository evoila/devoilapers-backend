package demolib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_HelloWorldEqualsHelloWorld(t *testing.T) {
	assert.Equal(t, "Hello World", "Hello World", "They should be equal")
}
