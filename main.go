package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error

	// wait for the database - shouldn't really be necessary except we're starting them together in compose
	time.Sleep(10 * time.Second)

	dsn := "host=postgres port=5432 user=postgres dbname=postgres sslmode=disable password=postgres"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&waterReading{})
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", handleReading)

	server := &http.Server{
		Addr:    ":8888",
		Handler: r,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Println("server starting")
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			log.Printf("error on ListenAndServe() %v", err)
		}
	}()

	<-done
	log.Println("server stopping")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("failed to shutdown gracefully %v", err)
	}
	log.Println("stopped gracefully")
}

func handleReading(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var reading waterReading
	err := decoder.Decode(&reading)
	if err != nil {
		log.Printf("uh oh: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("bad request"))
		if err != nil {
			log.Printf("error writing bad request bytes %v", err)
		}
		return
	}
	err = storeReading(reading)
	if err != nil {
		log.Printf("uh oh: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("internal server error"))
		if err != nil {
			log.Printf("error writing internal server error %v", err)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(reading)
	if err != nil {
		log.Printf("error writing json response %v", err)
	}
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
	TimeStamp   int    `json:"timeStamp" gorm:"not null"`
	Symbol      string `json:"symbol" gorm:"not null"`
	Volume      int    `json:"volume" gorm:"not null"`
	Temperature int    `json:"temperature" gorm:"not null"`
}
