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

// creates a new team member of a merchant
func CreateMember(w http.ResponseWriter, r *http.Request) {
	var member models.Member

	err := json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		log.Printf("Unable to decode the request body.  %v", err)
		json.NewEncoder(w).Encode(response{ID: -1, Message: err.Error()})
		return
	}

	memberID, err := insertMember(&member)
	if err != nil {
		msg := fmt.Sprintf("Unable to create a new member. %v", err)
		log.Println(msg)
		json.NewEncoder(w).Encode(response{ID: memberID, Message: msg})
		return
	}

	log.Printf("Member created successfully. memberID=%v, name=%v, email=%v, merchantID=%v", memberID, member.Name, member.Email, member.MerchantID)

	json.NewEncoder(w).Encode(response{ID: memberID, Message: "Member created successfully."})
}

// gets the member info
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
		log.Printf("Unable to get member, memberID=%v. %v", memberID, err)
		json.NewEncoder(w).Encode(response{ID: -1, Message: err.Error()})
		return
	}

	log.Printf("Get member: memberID=%v, name=%v, email=%v, merchantID=%v", memberID, member.Name, member.Email, member.MerchantID)

	json.NewEncoder(w).Encode(member)
}

// updates member info. the email is checked to avoid being duplicate to other members
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
		msg := fmt.Sprintf("Unable to update a new member. %v", err)
		log.Println(msg)
		json.NewEncoder(w).Encode(response{ID: int64(memberID), Message: msg})
		return
	}

	log.Printf("memberID=%v, %v\n", memberID, msg)

	json.NewEncoder(w).Encode(response{ID: int64(memberID), Message: msg})
}

// deletes a member by memberID
func DeleteMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	memberID, err := strconv.Atoi(params["memberid"])
	if err != nil {
		log.Printf("Unable to convert the string into int.  %v", err)
		json.NewEncoder(w).Encode(response{ID: -1, Message: err.Error()})
		return
	}

	msg := "Deleted a member successfully"
	err = deleteMember(int64(memberID))
	if err != nil {
		msg = "Unable to delete a member"
	}

	log.Printf("%v. memberID=%v\n", msg, memberID)

	json.NewEncoder(w).Encode(response{ID: int64(memberID), Message: msg})
}

//-- private methods

// connects to database, executes a query that inserts a new member. email is checked to avoid being duplicate beforhand
func insertMember(member *models.Member) (int64, error) {
	// check if email exists
	exist, err := checkIfEmailExists(-1, member.Email)
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

// connects to database, executes a select query to get member info by memberID
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

// connects to database, executes an update query on a member. email is checked to avoid being duplicate beforehand
func updateMember(memberID int64, member *models.Member) error {
	// check if the updated email is duplicate to others
	emailExist, err := checkIfEmailExists(memberID, member.Email)
	if err == nil && emailExist {
		return errors.New("email already exists")
	}

	// update data on db
	db := DbConnect()

	sqlStatement := `UPDATE members SET name=$2, email=$3, merchantID=$4 WHERE memberID=$1`

	_, err = db.Exec(sqlStatement, memberID, member.Name, member.Email, member.MerchantID)
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
	}

	return err
}

// connects to database, executes a delete query to delete a member by memberID
func deleteMember(memberID int64) error {
	db := DbConnect()

	sqlStatement := `DELETE FROM members WHERE memberID=$1`

	_, err := db.Exec(sqlStatement, memberID)
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
	}

	return err
}

// check if an email exists in members table. returns true with no error if it does, and false otherwise
func checkIfEmailExists(memberID int64, email string) (bool, error) {
	if email == "" {
		return false, errors.New("emails must be non-empty")
	}

	db := DbConnect()

	// check if email exists
	sqlStatement := `SELECT email FROM members WHERE email=$1 AND memberID<>$2`

	rows, err := db.Query(sqlStatement, email, memberID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}
