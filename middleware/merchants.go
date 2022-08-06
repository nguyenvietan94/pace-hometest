package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

	// TODO: lowercase

	// create an empty Merchant of type models.Merchant
	var Merchant models.Merchant

	// decode the json request to Merchant
	err := json.NewDecoder(r.Body).Decode(&Merchant)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert Merchant function and pass the Merchant
	insertID := insertMerchant(Merchant)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "Merchant created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func GetMerchant(w http.ResponseWriter, r *http.Request) {
	// get the Merchantid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the getMerchant function with Merchant id to retrieve a single Merchant
	Merchant, err := getMerchant(int64(id))

	if err != nil {
		log.Fatalf("Unable to get Merchant. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(Merchant)
}

// GetAllMerchant will return all the Merchants
func GetAllMerchants(w http.ResponseWriter, r *http.Request) {

	// get all the Merchants in the db
	Merchants, err := getAllMerchants()

	if err != nil {
		log.Fatalf("Unable to get all Merchant. %v", err)
	}

	// send all the Merchants as response
	json.NewEncoder(w).Encode(Merchants)
}

func UpdateMerchant(w http.ResponseWriter, r *http.Request) {

	// get the Merchantid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// create an empty Merchant of type models.Merchant
	var Merchant models.Merchant

	// decode the json request to Merchant
	err = json.NewDecoder(r.Body).Decode(&Merchant)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call update Merchant to update the Merchant
	updatedRows := updateMerchant(int64(id), Merchant)

	// format the message string
	msg := fmt.Sprintf("Merchant updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func DeleteMerchant(w http.ResponseWriter, r *http.Request) {

	// get the Merchantid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the deleteMerchant, convert the int to int64
	deletedRows := deleteMerchant(int64(id))

	// format the message string
	msg := fmt.Sprintf("Merchant updated successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------

// TODO: pass by pointer

// insert one Merchant in the DB
func insertMerchant(Merchant models.Merchant) int64 {

	// create the postgres db connection
	db := DbConnect()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning Merchantid will return the id of the inserted Merchant
	sqlStatement := `INSERT INTO merchants (name, age, location) VALUES ($1, $2, $3) RETURNING merchantID` // TODO: create table again

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, Merchant.Name, Merchant.Location, Merchant.Age).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

// get one Merchant from the DB by its Merchantid
func getMerchant(id int64) (models.Merchant, error) {
	// create the postgres db connection
	db := DbConnect()

	// close the db connection
	defer db.Close()

	// create a Merchant of models.Merchant type
	var Merchant models.Merchant

	// create the select sql query
	sqlStatement := `SELECT * FROM merchants WHERE merchantID=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to Merchant
	err := row.Scan(&Merchant.MerchantID, &Merchant.Name, &Merchant.Age, &Merchant.Location)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return Merchant, nil
	case nil:
		return Merchant, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty Merchant on error
	return Merchant, err
}

// get one Merchant from the DB by its Merchantid
func getAllMerchants() ([]models.Merchant, error) {
	// create the postgres db connection
	db := DbConnect()

	// close the db connection
	defer db.Close()

	var Merchants []models.Merchant

	// create the select sql query
	sqlStatement := `SELECT * FROM merchants`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var Merchant models.Merchant

		// unmarshal the row object to Merchant
		err = rows.Scan(&Merchant.MerchantID, &Merchant.Name, &Merchant.Age, &Merchant.Location)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the Merchant in the Merchants slice
		Merchants = append(Merchants, Merchant)

	}

	// return empty Merchant on error
	return Merchants, err
}

// update Merchant in the DB
func updateMerchant(id int64, Merchant models.Merchant) int64 {

	// create the postgres db connection
	db := DbConnect()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE merchants SET name=$2, location=$3, age=$4 WHERE merchantID=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, Merchant.Name, Merchant.Location, Merchant.Age)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete Merchant in the DB
func deleteMerchant(id int64) int64 {

	// create the postgres db connection
	db := DbConnect()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM merchants WHERE merchantID=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
