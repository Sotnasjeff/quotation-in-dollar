package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	CREATE_TABLE = "CREATE TABLE IF NOT EXISTS quotation(code TEXT, codein TEXT, name TEXT, high TEXT, low TEXT, varbid TEXT, pctchange TEXT, bid TEXT, ask TEXT, timestamp TEXT, created_at TEXT)"
	INSERT_INTO  = "INSERT INTO quotation(code, codein, name, high, low, varbid, pctchange, bid, ask, timestamp, created_at) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
)

type QuotationResponse struct {
	USDBRL USDBRL
}

type USDBRL struct {
	Code         string `json:"code"`
	Codein       string `json:"codein"`
	Name         string `json:"name"`
	High         string `json:"high"`
	Low          string `json:"low"`
	VarBid       string `json:"varBid"`
	PctChange    string `json:"pctChange"`
	Bid          string `json:"bid"`
	Ask          string `json:"ask"`
	Timestamp    string `json:"timestamp"`
	Created_date string `json:"create_date"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", HomeHandler)
	http.ListenAndServe(":8080", mux)

}

// "You See, Companies, They Come And Go. But Talent...Talent Is Forever" - Homelander
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	db, err := databaseConnection()
	if err != nil {
		log.Fatalf("Error in connecting to Database %v", err)
	}
	err = createTable(db)
	if err != nil {
		log.Fatalf("Error in creating table: %v", err)
	}
	defer db.Close()
	req, err := quotationDollarHttpClient(ctx)
	if err != nil {
		log.Fatalf("Error in making request to quotation api %v", err)
	}
	err = insertIntoDatabase(ctx, db, req)
	if err != nil {
		log.Fatalf("Error in inserting request to databse %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req.USDBRL.Bid)
}

func quotationDollarHttpClient(ctx context.Context) (*QuotationResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error in doing the request %v", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error in reading body %v", err)
		return nil, err
	}
	var quotation QuotationResponse
	err = json.Unmarshal(body, &quotation)
	if err != nil {
		log.Fatalf("Error in parsing response %v", err)
		return nil, err
	}

	return &quotation, nil

}

func databaseConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "quotation.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(CREATE_TABLE)
	return err
}

func insertIntoDatabase(ctx context.Context, db *sql.DB, quotation *QuotationResponse) error {
	stmt, err := db.Prepare(INSERT_INTO)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		quotation.USDBRL.Code,
		quotation.USDBRL.Codein,
		quotation.USDBRL.Name,
		quotation.USDBRL.High,
		quotation.USDBRL.Low,
		quotation.USDBRL.VarBid,
		quotation.USDBRL.PctChange,
		quotation.USDBRL.Bid,
		quotation.USDBRL.Ask,
		quotation.USDBRL.Timestamp,
		quotation.USDBRL.Created_date,
	)

	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		log.Fatalf("Timeout for inserting db")
		return err
	case <-time.After(10 * time.Nanosecond):
		return nil
	}
}
