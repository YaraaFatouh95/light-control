package controller

import (
	"encoding/json"
	"light-control/database"
	"light-control/models"

	"net/http"

	"github.com/gorilla/mux"
)

func CreateZone(w http.ResponseWriter, r *http.Request) {
	var newZone models.Zone
	err := json.NewDecoder(r.Body).Decode(&newZone)
	if err != nil {
		http.Error(w, "Invalid", http.StatusBadRequest)
		return
	}

	result := database.DB.Create(&newZone)
	if result.Error != nil {
		http.Error(w, "Failed to create zone", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(newZone)
}

func GetZones(w http.ResponseWriter, r *http.Request) {
	var zones []models.Zone
	result := database.DB.Preload("Luminaires").Find(&zones)
	if result.Error != nil {
		http.Error(w, "Failed to fetch zones", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zones)
}

func GetZone(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var zone models.Zone
	if err := database.DB.Preload("Luminaires").First(&zone, "id = ?", id).Error; err != nil {
		http.Error(w, "Zone not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zone)

}

func UpdateZone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var zone models.Zone

	if err := database.DB.First(&zone, "id = ?", id).Error; err != nil {
		http.Error(w, "Zone not found", http.StatusNotFound)
		return
	}

	var updatedData models.Zone
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	zone.Name = updatedData.Name
	zone.CityID = updatedData.CityID

	if err := database.DB.Save(zone).Error; err != nil {
		http.Error(w, "Failed to update zone", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zone)
}

func DeleteZone(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var luminaire models.Luminaire

	if err := database.DB.First(&luminaire, "zone_id = ?", id).Error; err == nil {
		http.Error(w, "Zone has luminaires and can not be deleted", http.StatusBadRequest)
		return
	}

	var zone models.Zone
	if err := database.DB.First(&zone, "id = ?", id).Error; err != nil {
		http.Error(w, "Zone not found", http.StatusNotFound)
		return
	}

	database.DB.Exec("DELETE FROM command_zones WHERE zone_id = ?", zone.ID)

	if err := database.DB.Delete(&zone).Error; err != nil {
		http.Error(w, "Failed to delete zone", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Zone deleted successfully"})
}
