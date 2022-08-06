package router

import (
	"pace-hometest/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// CRUD for a merchant account
	router.HandleFunc("/api/newmerchant", middleware.CreateMerchant).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/merchant/{id}", middleware.GetMerchant).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/merchant/{id}", middleware.UpdateMerchant).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletemerchant/{id}", middleware.DeleteMerchant).Methods("DELETE", "OPTIONS")

	// CRUD for a team member of a merchant account
	router.HandleFunc("/api/merchant/{merchantid}/newmember", middleware.CreateMember).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/merchant/{merchantid}/{memberid}", middleware.GetMember).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/merchant/{merchantid}/{memberid}", middleware.UpdateMember).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletemember/{merchantid}/{memberid}", middleware.DeleteMember).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/merchant/{merchantid}/allmembers", middleware.GetAllMembers).Methods("GET", "OPTIONS")

	return router
}
