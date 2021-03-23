package utils

import (
	"github.com/Pallinder/go-randomdata"
	"math/rand"
	"strings"
)

func GetRandomKubernetesResourceName(typePrefix string) string {
	prefix := typePrefix

	for i := 0; i < 2; i++ {
		min := 0
		max := 6
		value := rand.Intn(max-min) + min

		switch value {
		case 0:
			{
				prefix += "-" + randomdata.City()
			}
		case 1:
			{
				prefix += "-" + randomdata.Noun()
			}
		case 2:
			{
				prefix += "-" + randomdata.Adjective()
			}
		case 3:
			{
				prefix += "-" + randomdata.Country(randomdata.FullCountry)
			}
		case 4:
			{
				prefix += "-" + randomdata.FirstName(randomdata.RandomGender)
			}
		case 5:
			{
				prefix += "-" + randomdata.LastName()
			}
		}
	}
	// Remove spaces
	prefix = strings.Replace(prefix, " ", "", -1)
	return strings.ToLower(prefix)
}
