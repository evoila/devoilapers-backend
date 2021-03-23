package util

import (
	"OperatorAutomation/pkg/utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func Test_RandomValues(t *testing.T) {
	rand.Seed(42)
	knownValues := map[string]bool{}

	for i := 0; i < 300; i++ {
		randomName := utils.GetRandomKubernetesResourceName("test")
		_, exists := knownValues[randomName]

		if exists {
			fmt.Println(randomName)
		}

		assert.False(t, exists)
		knownValues[randomName] = true
	}
}
