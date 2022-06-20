package routes

import (
	"Toy_Cryptocurrency/controllers"

	"github.com/gorilla/mux"
)

func UserRoute(router *mux.Router) {
	router.HandleFunc("/sendSecurityCodeLogin", controllers.SendSecurityCodeLogin()).Methods("POST")
	router.HandleFunc("/verifySecurityCodeLogin/{securityCode}", controllers.VerifySecurityCodeLogin()).Methods("POST")
	router.HandleFunc("/sendSecurityCodeRegister", controllers.SendSecurityCodeRegister()).Methods("POST")
	router.HandleFunc("/verifySecurityCodeRegister/{securityCode}", controllers.VerifySecurityCodeRegister()).Methods("POST")
	router.HandleFunc("/getAvailableUsers/{userEmail}", controllers.GetAvailableUsers()).Methods("GET")
	router.HandleFunc("/getBalance/{userEmail}", controllers.GetBalance()).Methods("GET")
}
