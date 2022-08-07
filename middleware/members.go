package middleware

import (
	"encoding/json"
	"errors"
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
		log.Printf("Unable to decode the request body.  %v", err)
		json.NewEncoder(w).Encode(response{ID: -1, Message: err.Error()})
		return
	}

	msg := "Member created sucessfully."
	memberID, err := insertMember(&member)
	if err != nil {
		msg = "Unable to create a new member."
	}

	log.Printf("%v. memberID=%v", msg, memberID)

	json.NewEncoder(w).Encode(response{ID: memberID, Message: msg})
}

func GetMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	memberID, err := strconv.Atoi(params["memberid"])
	if err != nil {
		log.Printf("Unable to convert the string into int.  %v", err)
		json.NewEncoder(w).Encode(response{ID: -1, Message: err.Error()})
		return
	}

	member, err := getMember(int64(memberID))
	if err != nil {
		log.Printf("Unable to get member, ID=%v. %v", memberID, err)
	}

	log.Printf("Get member: memberID=%v, name=%v, email=%v, merchantID=%v", memberID, member.Name, member.Email, member.MerchantID)

	json.NewEncoder(w).Encode(member)
}

func UpdateMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	memberID, err := strconv.Atoi(params["memberid"])
	if err != nil {
		log.Printf("Unable to convert the string into int.  %v", err)
		json.NewEncoder(w).Encode(response{ID: -1, Message: err.Error()})
		return
	}

	var member models.Member
	err = json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		log.Printf("Unable to decode the request body.  %v", err)
		json.NewEncoder(w).Encode(response{ID: -1, Message: err.Error()})
		return
	}

	msg := "Member updated successfully."
	err = updateMember(int64(memberID), &member)
	if err != nil {
		msg = "Unable to update a member."
	}

	json.NewEncoder(w).Encode(response{ID: int64(memberID), Message: msg})
}

func DeleteMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	memberID, err := strconv.Atoi(params["memberid"])
	if err != nil {
		fmt.Printf("Unable to convert the string into int.  %v\n", err)
		// TODO: error should return empty message?
	}

	msg := "Deleted a member successfully."
	err = deleteMember(int64(memberID))
	if err != nil {
		msg = "Unable to delete a member."
	}
	fmt.Printf("memberID=%v, %v\n", memberID, msg)

	res := response{
		ID:      int64(memberID),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

//-- private methods ---

// TODO: return error, not return integer
func insertMember(member *models.Member) (int64, error) {
	// check if email exists
	exist, err := checkIfEmailExists(member.Email)
	if err == nil && exist {
		return -1, errors.New("email already exists")
	}

	// insert a new member to db
	db := DbConnect()

	sqlStatement := `INSERT INTO members (name, email, merchantID) VALUES ($1, $2, $3) RETURNING memberID`

	var memberID int64
	err = db.QueryRow(sqlStatement, member.Name, member.Email, member.MerchantID).Scan(&memberID)
	if err != nil {
		log.Printf("Unable to execute the query. %v\n", err)
		return -1, err
	}

	log.Printf("MerchantID %d: Inserted a new team member with memberID %d\n", member.MerchantID, memberID)

	return memberID, nil
}

func getMember(memberID int64) (*models.Member, error) {
	db := DbConnect()

	sqlStatement := `SELECT * FROM members WHERE memberID=$1`
	row := db.QueryRow(sqlStatement, memberID)

	var member models.Member
	err := row.Scan(&member.MemberID, &member.Name, &member.Email, &member.MerchantID)
	if err != nil {
		log.Printf("Unable to scan the row. %v", err)
		return nil, err
	}

	return &member, nil
}

func updateMember(memberID int64, member *models.Member) error {
	// check if email exists
	exist, err := checkIfEmailExists(member.Email)
	if err == nil && exist {
		return errors.New("email already exists")
	}

	// update data on db
	db := DbConnect()

	sqlStatement := `UPDATE members SET name=$2, email=$3, merchantID=$4 WHERE memberID=$1`

	_, err = db.Exec(sqlStatement, memberID, member.Name, member.Email, member.MerchantID)
	if err != nil {
		fmt.Printf("Unable to execute the query. %v\n", err)
	}

	return err
}

func deleteMember(memberID int64) error {
	db := DbConnect()

	sqlStatement := `DELETE FROM members WHERE memberID=$1`

	_, err := db.Exec(sqlStatement, memberID)
	if err != nil {
		fmt.Printf("Unable to execute the query. %v\n", err)
	}

	return err
}

func checkIfEmailExists(email string) (bool, error) {
	if email == "" {
		return false, errors.New("emails must be non-empty")
	}

	db := DbConnect()

	// check if email exists
	sqlStatement := `SELECT email FROM members WHERE email=$1`

	rows, err := db.Query(sqlStatement, email)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}
