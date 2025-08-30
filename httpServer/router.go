package httpServer

import (
	"light-control/controller"

	"github.com/gorilla/mux"
)

func RegisterCityRoutes(router *mux.Router) {
	router.HandleFunc("/cities", controller.CreateCity).Methods("POST")
	router.HandleFunc("/cities", controller.GetCities).Methods("GET")
	router.HandleFunc("/cities/{id}", controller.GetCity).Methods("GET")
	router.HandleFunc("/cities/{id}", controller.UpdateCity).Methods("PUT")
	router.HandleFunc("/cities/{id}", controller.DeleteCity).Methods("DELETE")
}

func RegisterZoneRoutes(router *mux.Router) {
	router.HandleFunc("/zones", controller.CreateZone).Methods("POST")
	router.HandleFunc("/zones", controller.GetZones).Methods("GET")
	router.HandleFunc("/zones/{id}", controller.GetZone).Methods("GET")
	router.HandleFunc("/zones/{id}", controller.UpdateZone).Methods("PUT")
	router.HandleFunc("/zones/{id}", controller.DeleteZone).Methods("DELETE")
}

func RegisterLuminaireRoutes(router *mux.Router) {
	router.HandleFunc("/luminaire", controller.CreateLuminaire).Methods("POST")
	router.HandleFunc("/luminaire", controller.GetLuminaires).Methods("GET")
	router.HandleFunc("/luminaire/{id}", controller.GetLuminaire).Methods("GET")
	router.HandleFunc("/luminaire/{id}", controller.UpdateLuminaire).Methods("PUT")
	router.HandleFunc("/luminaire/{id}", controller.DeleteLuminaire).Methods("DELETE")
}

func RegisterCommandRoutes(router *mux.Router) {
	router.HandleFunc("/command", controller.CreateCommand).Methods("POST")
	router.HandleFunc("/command", controller.GetCommands).Methods("GET")
	router.HandleFunc("/command/{id}", controller.GetCommand).Methods("GET")
	router.HandleFunc("/command/{id}", controller.UpdateCommand).Methods("PUT")
	router.HandleFunc("/command/{id}", controller.DeleteCommand).Methods("DELETE")
	router.HandleFunc("/command/exec", controller.ExecCommand).Methods("POST")
}
