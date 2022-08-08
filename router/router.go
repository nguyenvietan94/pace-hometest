package router

import (
	"pace-hometest/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// CRUD for merchant accounts
	router.HandleFunc("/api/newmerchant", middleware.CreateMerchant).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/merchant/{id}", middleware.GetMerchant).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/merchant/{id}", middleware.UpdateMerchant).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletemerchant/{id}", middleware.DeleteMerchant).Methods("DELETE", "OPTIONS")

	// CRUD for team members of a merchant account
	router.HandleFunc("/api/newmember", middleware.CreateMember).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/member/{memberid}", middleware.GetMember).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/member/{memberid}", middleware.UpdateMember).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletemember/{memberid}", middleware.DeleteMember).Methods("DELETE", "OPTIONS")

	// get team member lists with pagination
	router.HandleFunc("/api/merchant/{id}/allmembers", middleware.GetMembersWithPagination).Methods("GET", "OPTIONS")

	return router
}
