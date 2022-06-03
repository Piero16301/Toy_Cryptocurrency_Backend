package main

import (
	"Toy_Cryptocurrency/configs"
	"Toy_Cryptocurrency/routes"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	// Establecer conexión con MongoDB
	configs.ConnectDB()

	// Rutas de Blockchain
	routes.BlockRoute(router)

	// Rutas de Users
	routes.UserRoute(router)

	log.Fatal(http.ListenAndServe(":80", router))
}
