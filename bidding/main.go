package main

import (
	"bidding/controllers"
	"log"
	"net/http"
)

func main() {
	log.Println("Bidding Service is up and running...")
	http.HandleFunc("/bid", controllers.BidHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
