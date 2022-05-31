package routes

import (
	"Toy_Cryptocurrency/controllers"
	"github.com/gorilla/mux"
)

func BlockRoute(router *mux.Router) {
	router.HandleFunc("/newTransaction", controllers.NewTransaction()).Methods("POST")
	router.HandleFunc("/getBlockchain", controllers.GetBlockchain()).Methods("GET")
	router.HandleFunc("/validateBlockchain", controllers.ValidateBlockchain()).Methods("GET")
}
