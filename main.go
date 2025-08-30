package main

import (
	"light-control/database"
	"light-control/httpServer"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	//"light-control/routes"
)

func main() {
	database.ConnectDB()

	r := mux.NewRouter()
	httpServer.RegisterCityRoutes(r)
	httpServer.RegisterZoneRoutes(r)
	httpServer.RegisterLuminaireRoutes(r)
	httpServer.RegisterCommandRoutes(r)

	log.Println("Server running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
