package utils

import (
	"OperatorAutomation/pkg/core/common"
	"github.com/Pallinder/go-randomdata"
	"strings"
)

func FillWithData(authInformation common.IKubernetesAuthInformation, template string) string {
	template = strings.Replace(template, "{opa.random.name}", getRandomName(), -1)
	template = strings.Replace(template, "{opa.user.namespace}", authInformation.GetKubernetesNamespace(), -1)

	return template
}

func getRandomName() string {
	randomName := randomdata.FirstName(randomdata.RandomGender)
	if randomdata.Boolean() {
		randomName = randomName + "-" + randomdata.City()
	} else {
		randomName = randomdata.City() + "-" + randomName
	}

	// Remove spaces
	randomName = strings.Replace(randomName, " ", "", -1)
	return strings.ToLower(randomName)
}