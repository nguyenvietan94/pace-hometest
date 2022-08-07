package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"pace-hometest/models"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateMember(w http.ResponseWriter, r *http.Request) {
	var member models.Member

	err := json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		fmt.Printf("Unable to decode the request body.  %v", err)
		return
	}

	msg := "Member created sucessfully."
	memberID := insertMember(&member)
	if memberID < 0 {
		msg = "Unable to create a new member."
	}

	res := response{
		ID:      memberID,
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func GetMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	memberID, err := strconv.Atoi(params["memberid"])
	if err != nil {
		fmt.Printf("Unable to convert the string into int.  %v\n", err)
	}

	member, err := getMember(int64(memberID))
	if err != nil {
		fmt.Printf("Unable to get a member. %v\n", err)
	}

	json.NewEncoder(w).Encode(member)
}

func UpdateMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	memberID, err := strconv.Atoi(params["memberid"])
	if err != nil {
		fmt.Printf("Unable to convert the string into int.  %v\n", err)
	}

	var member models.Member
	err = json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		fmt.Printf("Unable to decode the request body.  %v", err)
	}

	msg := "Member updated successfully."
	err = updateMember(int64(memberID), &member)
	if err != nil {
		msg = "Unable to update a member."
	}

	res := response{
		ID:      int64(memberID),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
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
func insertMember(member *models.Member) int64 {
	// check if email exists
	exist, err := checkIfEmailExists(member.Email)
	if err == nil && exist {
		return -1
	}

	// insert a new member to db
	db := DbConnect()

	sqlStatement := `INSERT INTO members (name, email, merchantID) VALUES ($1, $2, $3) RETURNING memberID`

	var memberID int64
	err = db.QueryRow(sqlStatement, member.Name, member.Email, member.MerchantID).Scan(&memberID)
	if err != nil {
		fmt.Printf("Unable to execute the query. %v\n", err)
		return -1
	}

	fmt.Printf("MerchantID %d: Inserted a new team member with memberID %d\n", member.MerchantID, memberID)

	return memberID
}

func getMember(memberID int64) (*models.Member, error) {
	db := DbConnect()

	sqlStatement := `SELECT * FROM members WHERE memberID=$1`
	row := db.QueryRow(sqlStatement, memberID)

	var member models.Member
	err := row.Scan(&member.MemberID, &member.Name, &member.Email, &member.MerchantID)
	if err != nil {
		fmt.Printf("Unable to scan the row. %v\n", err)
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
