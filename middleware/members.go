package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pace-hometest/models"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateMember(w http.ResponseWriter, r *http.Request) {
	var member models.Member

	err := json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	memId := insertMember(&member)

	res := response{
		ID:      memId,
		Message: "Member created sucessfully",
	}

	json.NewEncoder(w).Encode(res)
}

func GetMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	memberId, err := strconv.Atoi(params["memberid"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	member, err := getMember(int64(memberId))
	if err != nil {
		log.Fatalf("Unable to get a member. %v", err)
	}

	json.NewEncoder(w).Encode(member)
}

func UpdateMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	memberId, err := strconv.Atoi(params["memberid"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	var member models.Member
	err = json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	var msg string
	err = updateMember(int64(memberId), &member)
	if err == nil {
		msg = fmt.Sprintf("Updating member: succeeded.")
	} else {
		msg = fmt.Sprintf("Updating member: failed.")
	}

	res := response{
		ID:      int64(member.Id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	memberId, err := strconv.Atoi(params["memberid"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	var msg string
	err = deleteMember(int64(memberId))
	if err == nil {
		msg = fmt.Sprintf("Deleting member: succeeded.")
	} else {
		msg = fmt.Sprintf("Deleting member: failed.")
	}

	res := response{
		ID:      int64(memberId),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// TODO: implement pagination
func GetAllMembers(w http.ResponseWriter, r *http.Request) {

}

//-- private methods

func insertMember(member *models.Member) int64 {
	db := DbConnect()
	defer db.Close()

	sqlStatement := `INSERT INTO members (merchantId, name, email) VALUES ($1, $2, $3) RETURNING id`

	var memberId int64
	err := db.QueryRow(sqlStatement, member.MerchantId, member.Name, member.Email).Scan(&memberId)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Merchant %d: Inserted a new team member with Id %d\n", member.MerchantId, memberId)

	return memberId
}

func getMember(memberId int64) (*models.Member, error) {
	db := DbConnect()
	defer db.Close()

	sqlStatement := `SELECT * FROM members WHERE id=$1`
	row := db.QueryRow(sqlStatement, memberId)

	var member models.Member
	err := row.Scan(&member.Id, &member.MerchantId, &member.Name, &member.Email)
	if err != nil {
		log.Fatalf("Unable to scan the row. %v", err)
		return nil, err
	}

	return &member, nil
}

func updateMember(memberId int64, member *models.Member) error {
	db := DbConnect()
	defer db.Close()

	sqlStatement := `UPDATE Merchants SET merchantId=$1, name=$2, email=$3 WHERE id=$4`

	_, err := db.Exec(sqlStatement, member.MerchantId, member.Name, member.Email, memberId)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	return err
}

func deleteMember(memberId int64) error {
	db := DbConnect()
	defer db.Close()

	sqlStatement := `DELETE FROM members WHERE id=$1`

	_, err := db.Exec(sqlStatement, memberId)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	return err
}
