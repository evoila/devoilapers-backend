package utils

import (
	"github.com/Pallinder/go-randomdata"
	"strings"
)

func GetRandomKubernetesResourceName() string {
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
