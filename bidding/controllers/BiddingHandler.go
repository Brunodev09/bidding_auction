package controllers

import (
	internal "bidding/internal"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
)

type AdObject struct {
	AdID     uuid.UUID `json:"ad_id"`
	BidPrice int       `json:"bid_price"`
}

func BidHandler(w http.ResponseWriter, r *http.Request) {

	// Simulate not bidding and returning no content
	if rand.Intn(10) < 2 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Simulate a random bid price
	bidPrice := rand.Intn(100)
	adObject := AdObject{
		AdID:     uuid.New(),
		BidPrice: bidPrice,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(adObject)

	adObjectBytes, err := json.Marshal(adObject)

	if err != nil {
		log.Fatal(err)
	}

	internal.PushEventToQueue("bids", adObjectBytes)
}
