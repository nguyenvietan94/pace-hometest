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
	router.HandleFunc("/api/merchant/{id}/allmembers", middleware.GetMembersWithPagination).Methods("GET", "OPTIONS")

	// CRUD for a team member of a merchant account
	router.HandleFunc("/api/member/newmember", middleware.CreateMember).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/member/{memberid}", middleware.GetMember).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/member/{memberid}", middleware.UpdateMember).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletemember/{memberid}", middleware.DeleteMember).Methods("DELETE", "OPTIONS")

	return router
}
