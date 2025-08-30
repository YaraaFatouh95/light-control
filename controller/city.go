package controller

import (
	"encoding/json"
	"light-control/database"
	"light-control/models"

	"net/http"

	"github.com/gorilla/mux"
)

func CreateCity(w http.ResponseWriter, r *http.Request) {
	var newCity models.City
	err := json.NewDecoder(r.Body).Decode(&newCity)
	if err != nil {
		http.Error(w, "Invalid", http.StatusBadRequest)
		return
	}

	result := database.DB.Create(&newCity)
	if result.Error != nil {
		http.Error(w, "Failed to create city", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newCity)
}

func GetCities(w http.ResponseWriter, r *http.Request) {
	var cities []models.City
	if err := database.DB.Preload("Zones").Find(&cities).Error; err != nil {
		http.Error(w, "Failed to fetch cities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cities)
}

func GetCity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var city models.City
	if err := database.DB.Preload("Zones").First(&city, "id = ?", id).Error; err != nil {
		http.Error(w, "City not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(city)
}

func UpdateCity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var city models.City
	err := database.DB.First(&city, "id = ?", id).Error
	if err != nil {
		http.Error(w, "City not found", http.StatusNotFound)
		return
	}

	var updatedData models.City
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	city.Name = updatedData.Name
	if err := database.DB.Save(&city).Error; err != nil {
		http.Error(w, "Failed to update city", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(city)
}
func DeleteCity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var zone models.Zone

	if err := database.DB.First(&zone, "city_id = ?", id).Error; err == nil {
		http.Error(w, "City has zones and can not be deleted", http.StatusBadRequest)
		return
	}

	var city models.City
	if err := database.DB.First(&city, "id = ?", id).Error; err != nil {
		http.Error(w, "City not found", http.StatusNotFound)
		return
	}

	database.DB.Exec("DELETE FROM command_cities WHERE city_id = ?", city.ID)

	if err := database.DB.Delete(&models.City{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Failed to delete city", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "City deleted successfully"})
}
