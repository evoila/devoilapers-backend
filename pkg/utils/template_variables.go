package utils

import (
	"OperatorAutomation/pkg/core/common"
	"github.com/Pallinder/go-randomdata"
	"regexp"
	"strings"
)

func FillWithData(authInformation common.IKubernetesAuthInformation, template string) string {
	template = strings.Replace(template, "{opa.random.name}", getRandomName(), -1)
	template = strings.Replace(template, "{opa.user.namespace}", authInformation.GetKubernetesNamespace(), -1)


	// Replace random names with index
	randomNameRegex := "{(opa.random.name)\\[(\\d)\\]}"
	regexPattern := regexp.MustCompile(randomNameRegex)
	matches := regexPattern.FindAllStringSubmatch(template, -1)
	// Indicates whenever a random index allready got replaced in the whole text
	allReadyReplaced := map[string]bool{}
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		wholeMatch := match[0]
		index := match[2]

		_, exists := allReadyReplaced[index]
		if !exists {
			template = strings.Replace(template, wholeMatch, getRandomName(), -1)
			allReadyReplaced[index] = true
		}
	}

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