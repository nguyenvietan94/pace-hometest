package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"pace-hometest/models"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateMerchant(w http.ResponseWriter, r *http.Request) {
	var merchant models.Merchant

	err := json.NewDecoder(r.Body).Decode(&merchant)
	if err != nil {
		fmt.Printf("Unable to decode the request body.  %v", err)
		return
	}

	msg := "Merchant created successfully"
	insertedID := insertMerchant(&merchant)
	if insertedID < 0 {
		msg = "Unable to create a merchant."
	}

	res := response{
		ID:      insertedID,
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func GetMerchant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Printf("Unable to convert the string into int.  %v\n", err)
		return
	}

	merchant, err := getMerchant(int64(id))

	if err != nil {
		fmt.Printf("Unable to get merchant. %v\n", err)
		return
	}

	json.NewEncoder(w).Encode(merchant)
}

func UpdateMerchant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Printf("Unable to convert the string into int.  %v", err)
		return
	}

	var merchant models.Merchant
	err = json.NewDecoder(r.Body).Decode(&merchant)
	if err != nil {
		fmt.Printf("Unable to decode the request body.  %v", err)
		return
	}

	msg := "Updated successfully"
	err = updateMerchant(int64(id), &merchant)
	if err != nil {
		msg = "Unable to update."
	}

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteMerchant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Printf("Unable to convert the string into int.  %v\n", err)
		return
	}

	msg := "Deleted successfully."
	err = deleteMerchant(int64(id))
	if err != nil {
		msg = "Unable to delete."
	}

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------

func insertMerchant(merchant *models.Merchant) int64 {
	db := DbConnect()

	sqlStatement := `INSERT INTO merchants (name, age, location) VALUES ($1, $2, $3) RETURNING merchantID` // TODO: create table again

	var id int64
	err := db.QueryRow(sqlStatement, merchant.Name, merchant.Age, merchant.Location).Scan(&id)
	if err != nil {
		fmt.Printf("Unable to execute the query. %v", err)
		return -1
	}

	fmt.Printf("Inserted a new merchant: id=%v, name=%v, age=%v, location=%v\n", id, merchant.Name, merchant.Age, merchant.Location)

	return id
}

func getMerchant(id int64) (*models.Merchant, error) {
	db := DbConnect()

	var merchant models.Merchant

	sqlStatement := `SELECT * FROM merchants WHERE merchantID=$1`
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&merchant.MerchantID, &merchant.Name, &merchant.Age, &merchant.Location)
	if err != nil {
		fmt.Printf("Unable to scan the row. %v\n", err)
		return nil, err
	}

	return &merchant, err
}

func updateMerchant(id int64, merchant *models.Merchant) error {
	db := DbConnect()

	sqlStatement := `UPDATE merchants SET name=$2, age=$3, location=$4 WHERE merchantID=$1`
	_, err := db.Exec(sqlStatement, id, merchant.Name, merchant.Age, merchant.Location)
	if err != nil {
		fmt.Printf("Unable to execute the query. %v\n", err)
	}

	return err
}

func deleteMerchant(id int64) error {
	db := DbConnect()

	sqlStatement := `DELETE FROM merchants WHERE merchantID=$1`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		fmt.Printf("Unable to execute the query. %v\n", err)
	}

	return err
}
