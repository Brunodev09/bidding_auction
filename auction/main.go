package main

import (
	"auction/controllers"
	"auction/internal"
	"log"
	"net/http"
)

func main() {
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

	http.HandleFunc("/auction", controllers.AuctionHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
