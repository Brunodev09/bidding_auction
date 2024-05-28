package main

import (
	"auction/controllers"
	"auction/internal"
	"auction/models"
	"auction/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_DATABASE"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal(err)
	}

	err = models.MigrateEvents(db)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server up and running on port 8080")

	// Start the Kafka consumer in a separate goroutine
	go func() {
		log.Println("Attempting to connect to Kafka broker...")
		if err := internal.EnableQueue(); err != nil {
			log.Fatalf("Failed to enable queue: %v", err)
		}
	}()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Auction Service, please use /auction endpoint to get the best bid."))
	})

	r := controllers.Repository{DB: db}

	http.HandleFunc("/auction", errorHandler(r.AuctionHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func errorHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			log.Printf("HTTP %d - %s", http.StatusInternalServerError, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
