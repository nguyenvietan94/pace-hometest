package middleware

import (
	"encoding/json"
	"errors"
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
	var merchant models.Merchant

	err := json.NewDecoder(r.Body).Decode(&merchant)
	if err != nil {
		log.Printf("Unable to decode the request body.  %v", err)
		json.NewEncoder(w).Encode(response{ID: -1, Message: err.Error()})
		return
	}

	msg := "Merchant created successfully"
	merchantID, err := insertMerchant(&merchant)
	if err != nil {
		msg = "Unable to create a merchant."
	}

	log.Printf("Inserted a new merchant: id=%v, name=%v, age=%v, location=%v\n", merchantID, merchant.Name, merchant.Age, merchant.Location)

	json.NewEncoder(w).Encode(response{ID: merchantID, Message: msg})
}

func GetMerchant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	merchantID, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Printf("Unable to convert the string into int.  %v\n", err)
		json.NewEncoder(w).Encode(response{ID: -1, Message: err.Error()})
		return
	}

	merchant, err := getMerchant(int64(merchantID))

	if err != nil {
		log.Printf("Unable to get merchant. %v\n", err)
		json.NewEncoder(w).Encode(response{ID: -1, Message: err.Error()})
		return
	}

	log.Printf("Get merchant: merchantID=%v, name=%v, email=%v, merchantID=%v", merchantID, merchant.Name, merchant.Age, merchant.Location)

	json.NewEncoder(w).Encode(merchant)
}

func UpdateMerchant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	merchantID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Printf("Unable to convert the string into int.  %v", err)
		json.NewEncoder(w).Encode(response{ID: -1, Message: err.Error()})
		return
	}

	var merchant models.Merchant
	err = json.NewDecoder(r.Body).Decode(&merchant)
	if err != nil {
		log.Printf("Unable to decode the request body.  %v", err)
		json.NewEncoder(w).Encode(response{ID: -1, Message: err.Error()})
		return
	}

	msg := "Updated successfully"
	err = updateMerchant(int64(merchantID), &merchant)
	if err != nil {
		msg = "Unable to update."
	}

	log.Printf("merchantID=%v, %v", merchantID, msg)

	json.NewEncoder(w).Encode(response{ID: int64(merchantID), Message: msg})
}

func DeleteMerchant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	merchantID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Printf("Unable to convert the string into int.  %v", err)
		json.NewEncoder(w).Encode(response{ID: -1, Message: err.Error()})
		return
	}

	msg := "Deleted a merchant successfully"
	err = deleteMerchant(int64(merchantID))
	if err != nil {
		msg = "Unable to delete a merchant"
	}

	log.Printf("%v, merchantID %v", msg, merchantID)

	json.NewEncoder(w).Encode(response{ID: int64(merchantID), Message: msg})
}

// implement pagination
func GetMembersWithPagination(w http.ResponseWriter, r *http.Request) {
	// get merchantID
	params := mux.Vars(r)

	merchantID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Printf("Unable to convert the string into int.  %v\n", err)
	}

	// get page id and page size
	pageID, err := strconv.Atoi(r.URL.Query().Get("pageid"))
	if err != nil {
		pageID = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pagesize"))
	if err != nil {
		pageSize = 10
	}

	members, err := getMembersWithPagination(int64(merchantID), pageID, pageSize)
	if err != nil {
		log.Printf("Unable to get members, merchantID=%v. %v\n", merchantID, err)
		json.NewEncoder(w).Encode(response{ID: -1, Message: err.Error()})
		return
	}

	log.Printf("Get team members: merchantID=%v, pageID=%v, pageSize=%v", merchantID, pageID, pageSize)

	json.NewEncoder(w).Encode(members)
}

//-- private methods

func insertMerchant(merchant *models.Merchant) (int64, error) {
	db := DbConnect()

	sqlStatement := `INSERT INTO merchants (name, age, location) VALUES ($1, $2, $3) RETURNING merchantID`

	var memberID int64
	err := db.QueryRow(sqlStatement, merchant.Name, merchant.Age, merchant.Location).Scan(&memberID)
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return -1, err
	}

	return memberID, nil
}

func getMerchant(id int64) (*models.Merchant, error) {
	db := DbConnect()

	var merchant models.Merchant

	sqlStatement := `SELECT * FROM merchants WHERE merchantID=$1`
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&merchant.MerchantID, &merchant.Name, &merchant.Age, &merchant.Location)
	if err != nil {
		log.Printf("Unable to scan the row. %v", err)
		return nil, err
	}

	return &merchant, err
}

func updateMerchant(id int64, merchant *models.Merchant) error {
	db := DbConnect()

	sqlStatement := `UPDATE merchants SET name=$2, age=$3, location=$4 WHERE merchantID=$1`
	_, err := db.Exec(sqlStatement, id, merchant.Name, merchant.Age, merchant.Location)
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
	}

	return err
}

func deleteMerchant(id int64) error {
	db := DbConnect()

	sqlStatement := `DELETE FROM merchants WHERE merchantID=$1`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
	}

	return err
}

// pageID >= 1
func getMembersWithPagination(merchantID int64, pageID, pageSize int) ([]models.Member, error) {
	if pageID <= 0 {
		return nil, errors.New("pageID must be greater than 0")
	}

	db := DbConnect()

	limit := pageID * pageSize
	start := (pageID - 1) * pageSize

	sqlStatement := `SELECT * FROM members WHERE merchantID=$1 LIMIT $2`
	rows, err := db.Query(sqlStatement, merchantID, limit)
	if err != nil {
		log.Printf("Unable to execute the query. %v\n", err)
		return nil, err
	}

	var members []models.Member
	cnt := 0
	for rows.Next() {
		if cnt >= start {
			var mem models.Member
			err = rows.Scan(&mem.MemberID, &mem.Name, &mem.Email, &mem.MerchantID)
			if err != nil {
				log.Printf("Unable to scan the row. %v\n", err)
			}
			members = append(members, mem)
		}
		cnt++
	}

	return members, nil
}
