package router

import (
	"pace-hometest/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// CRUD for a merchant account
	router.HandleFunc("/api/newuser", middleware.CreateMerchant).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user/{id}", middleware.GetMerchant).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user/{id}", middleware.UpdateMerchant).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deleteuser/{id}", middleware.DeleteMerchant).Methods("DELETE", "OPTIONS")

	// CRUD for a team member of a merchant account
	router.HandleFunc("/api/user/{id}/newmember", middleware.CreateMember).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user/{id}/{memberid}", middleware.GetMember).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user/{id}/{memberid}", middleware.UpdateMember).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletemember/{id}/{memberid}", middleware.DeleteMember).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/user/{id}/allmembers", middleware.GetAllMembers).Methods("GET", "OPTIONS")

	return router
}
