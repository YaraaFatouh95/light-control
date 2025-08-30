package utils

import (
	"fmt"
	"light-control/database"
	"light-control/models"

	"github.com/google/uuid"
)

func GenerateTopicToBeSent(entityType string, entities []uuid.UUID) ([]string, error) {
	var topics []string
	if entityType == "city" {
		var cities []models.City

		if err := database.DB.Preload("Zones.Luminaires").Find(&cities, "id IN ?", entities).Error; err != nil {
			return []string{}, fmt.Errorf("cities can not be fetched")
		}

		for _, city := range cities {
			for _, zone := range city.Zones {
				for _, luminaire := range zone.Luminaires {
					topics = append(topics, fmt.Sprintf("city/%v/zone/%v/luminaire/%v/command", city.Name, zone.Name, luminaire.Name))
				}
			}
		}
	} else if entityType == "zone" {
		var zones []models.Zone

		if err := database.DB.Preload("Luminaires").Find(&zones, "id IN ?", entities).Error; err != nil {
			return []string{}, fmt.Errorf("zones can not be fetched")
		}

		for _, zone := range zones {
			var city models.City
			if err := database.DB.First(&city, "id = ?", zone.CityID).Error; err != nil {
				return []string{}, fmt.Errorf("zone city can not be fetched")
			}
			for _, luminaire := range zone.Luminaires {

				topics = append(topics, fmt.Sprintf("city/%v/zone/%v/luminaire/%v/command", city.Name, zone.Name, luminaire.Name))
			}
		}
	} else {
		var luminaires []models.Luminaire

		if err := database.DB.Find(&luminaires, "id IN ?", entities).Error; err != nil {
			return []string{}, fmt.Errorf("luminaires can not be fetched")
		}

		for _, luminaire := range luminaires {
			var zone models.Zone
			if err := database.DB.First(&zone, "id = ?", luminaire.ZoneID).Error; err != nil {
				return []string{}, fmt.Errorf("zone city can not be fetched")
			}
			var city models.City
			if err := database.DB.First(&city, "id = ?", zone.CityID).Error; err != nil {
				return []string{}, fmt.Errorf("zone city can not be fetched")
			}
			topics = append(topics, fmt.Sprintf("city/%v/zone/%v/luminaire/%v/command", city.Name, zone.Name, luminaire.Name))
		}

	}

	return topics, nil
}
