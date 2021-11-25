package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

// todo initialized without global var
var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=water sslmode=disable password=postgres")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&waterReading{})

	// todo use a router like chi/gin to handle advanced routing
	// and add composablity with middlewares
	// this would take care of things like 405s
	http.HandleFunc("/", handleReading)
	log.Fatal(http.ListenAndServe(":8888", nil))
	// todo graceful shutdown
}

func handleReading(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var reading waterReading
	err := decoder.Decode(&reading)
	if err != nil {
		// todo something better with errors, like return to the client
		log.Printf("uh oh: %v\n", err)
	}
	err = storeReading(reading)
	w.WriteHeader(http.StatusCreated)
	return
}

func storeReading(reading waterReading) error {
	result := db.Create(&reading)
	if result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

type waterReading struct {
	gorm.Model
	TimeStamp   int    `json:"timeStamp"`
	Symbol      string `json:"symbol"`
	Volume      int    `json:"volume"`
	Temperature int    `json:"temperature"`
}
