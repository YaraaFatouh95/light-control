package controller

import (
	"encoding/json"
	"fmt"
	"light-control/database"
	"light-control/models"
	"light-control/utils"
	"time"

	"light-control/dkron"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
)

func CreateCommand(w http.ResponseWriter, r *http.Request) {
	var newCommand models.Command
	err := json.NewDecoder(r.Body).Decode(&newCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(newCommand.Entities) == 0 {
		http.Error(w, "entities can not be empty", http.StatusBadRequest)
		return
	}

	if newCommand.EntityType == "city" {
		var cities []models.City
		result := database.DB.Where("id IN ?", newCommand.Entities).Find(&cities)
		if result.Error != nil {
			http.Error(w, "Failed to fetch Cities", http.StatusInternalServerError)
			return
		}

		if len(newCommand.Entities) != len(cities) {
			http.Error(w, "Not All Citites are found", http.StatusBadRequest)
			return
		}
		newCommand.Cities = cities

	} else if newCommand.EntityType == "zone" {

		var zones []models.Zone
		result := database.DB.Where("id IN ?", newCommand.Entities).Find(&zones)
		if result.Error != nil {
			http.Error(w, "Failed to fetch Zones", http.StatusInternalServerError)
			return
		}

		if len(newCommand.Entities) != len(zones) {
			http.Error(w, "Not All Zones are found", http.StatusBadRequest)
			return
		}

		newCommand.Zones = zones

	} else if newCommand.EntityType == "luminaire" {
		var luminaires []models.Luminaire
		result := database.DB.Where("id IN ?", newCommand.Entities).Find(&luminaires)
		if result.Error != nil {
			http.Error(w, "Failed to fetch Luminaires", http.StatusInternalServerError)
			return
		}

		if len(newCommand.Entities) != len(luminaires) {
			http.Error(w, "Not All Luminaires are found", http.StatusBadRequest)
			return
		}

		newCommand.Luminaires = luminaires
	} else {
		http.Error(w, "entity_type should be city or zone or luminaire", http.StatusBadRequest)
		return
	}

	newCommand.Status = "Pending"
	if newCommand.ScheduledTime.Before(time.Now()) {
		http.Error(w, "Schedule Time can not be in the past", http.StatusBadRequest)
		return
	}

	result := database.DB.Create(&newCommand)
	if result.Error != nil {
		http.Error(w, "Failed to create command", http.StatusInternalServerError)
		return
	}
	if err := dkron.CreateDkronJob(newCommand); err != nil {
		http.Error(w, "Failed to create dkron job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newCommand)
}

func GetCommands(w http.ResponseWriter, r *http.Request) {
	var command []models.Command
	result := database.DB.Preload("Cities").Preload("Luminaires").Preload("Zones").Find(&command)
	if result.Error != nil {
		http.Error(w, "Failed to fetch Commands", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(command)
}

func GetCommand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var command models.Command
	if err := database.DB.Preload("Cities").Preload("Luminaires").Preload("Zones").First(&command, "id = ?", id).Error; err != nil {
		http.Error(w, "command not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(command)

}

func UpdateCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var command models.Command

	if err := database.DB.First(&command, "id = ?", id).Error; err != nil {
		http.Error(w, "command not found", http.StatusNotFound)
		return
	}

	var updatedData models.Command
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if updatedData.ScheduledTime.Before(time.Now()) {
		http.Error(w, "Schedule Time can not be in the past", http.StatusBadRequest)
		return
	}
	if command.ScheduledTime != updatedData.ScheduledTime {
		if err := dkron.DeleteDkronJob(command); err != nil {
			http.Error(w, "Failed to update dkron job: "+err.Error(), http.StatusInternalServerError)
			return
		}
		command.ScheduledTime = updatedData.ScheduledTime

		if err := dkron.CreateDkronJob(command); err != nil {
			http.Error(w, "Failed to update dkron job: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := database.DB.Save(command).Error; err != nil {
		http.Error(w, "Failed to update command", http.StatusInternalServerError)
		return
	}

	if err := database.DB.Preload("Cities").Preload("Luminaires").Preload("Zones").First(&command, "id = ?", id).Error; err != nil {
		http.Error(w, "Failed to fetch updated command", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(command)
}

func DeleteCommand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var command models.Command
	if err := database.DB.First(&command, "id = ?", id).Error; err != nil {
		http.Error(w, "command not found", http.StatusNotFound)
		return
	}

	if err := dkron.DeleteDkronJob(command); err != nil {
		http.Error(w, "Failed to delete dkron job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	database.DB.Exec("DELETE FROM command_cities WHERE command_id = ?", command.ID)
	database.DB.Exec("DELETE FROM command_zones WHERE command_id = ?", command.ID)
	database.DB.Exec("DELETE FROM command_luminaires WHERE command_id = ?", command.ID)

	if err := database.DB.Delete(&command).Error; err != nil {
		http.Error(w, "Failed to delete command", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "command deleted successfully"})
}

func ExecCommand(w http.ResponseWriter, r *http.Request) {

	var command models.Command
	err := json.NewDecoder(r.Body).Decode(&command)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := mqtt.NewClient(mqtt.NewClientOptions().
		AddBroker("tcp://broker.hivemq.com:1883").
		SetClientID("light-control-go-service"))

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		http.Error(w, token.Error().Error(), http.StatusBadRequest)
		return
	}

	topics, err := utils.GenerateTopicToBeSent(command.EntityType, command.Entities)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, topic := range topics {
		fmt.Println(topic)
		token := client.Publish(topic, 1, false, command.Payload)
		token.Wait()
	}

	client.Disconnect(250)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "command deleted successfully"})
}
