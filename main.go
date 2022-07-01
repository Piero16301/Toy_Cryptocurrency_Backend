package main

import (
	"Toy_Cryptocurrency/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Rutas de Blockchain
	routes.BlockRoute(router)

	// Rutas de Users
	routes.UserRoute(router)

	log.Fatal(http.ListenAndServe(":80", router))
}
