package controller

import (
	"encoding/json"
	"light-control/database"
	"light-control/models"

	"net/http"

	"github.com/gorilla/mux"
)

func CreateLuminaire(w http.ResponseWriter, r *http.Request) {
	var newLuminaire models.Luminaire
	err := json.NewDecoder(r.Body).Decode(&newLuminaire)
	if err != nil {
		http.Error(w, "Invalid"+err.Error(), http.StatusBadRequest)
		return
	}

	result := database.DB.Create(&newLuminaire)
	if result.Error != nil {
		http.Error(w, "Failed to create luminaire"+result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newLuminaire)
}

func GetLuminaires(w http.ResponseWriter, r *http.Request) {
	var luminaire []models.Luminaire
	result := database.DB.Find(&luminaire)
	if result.Error != nil {
		http.Error(w, "Failed to fetch luminaire"+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(luminaire)
}

func GetLuminaire(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var luminaire models.Luminaire
	if err := database.DB.First(&luminaire, "id = ?", id).Error; err != nil {
		http.Error(w, "luminaire not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(luminaire)

}

func UpdateLuminaire(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var luminaire models.Luminaire

	if err := database.DB.First(&luminaire, "id = ?", id).Error; err != nil {
		http.Error(w, "luminaire not found", http.StatusNotFound)
		return
	}

	var updatedData models.Luminaire
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	luminaire.Name = updatedData.Name
	luminaire.ZoneID = updatedData.ZoneID
	luminaire.Status = updatedData.Status
	luminaire.Dim = updatedData.Dim

	if err := database.DB.Save(luminaire).Error; err != nil {
		http.Error(w, "Failed to update luminaire", http.StatusInternalServerError)
		return
	}

	if err := database.DB.First(&luminaire, "id = ?", id).Error; err != nil {
		http.Error(w, "Failed to fetch updated luminaire", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(luminaire)
}

func DeleteLuminaire(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var luminaire models.Luminaire
	if err := database.DB.First(&luminaire, "id = ?", id).Error; err != nil {
		http.Error(w, "luminaire not found", http.StatusNotFound)
		return
	}

	database.DB.Exec("DELETE FROM command_luminaires WHERE luminaire_id = ?", luminaire.ID)

	if err := database.DB.Delete(&luminaire).Error; err != nil {
		http.Error(w, "Failed to delete luminaire", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "luminaire deleted successfully"})
}
