package controllers

import (
	"auction/models"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"gorm.io/gorm"
)

type BiddingServiceResponse struct {
	AdID     string `json:"ad_id"`
	BidPrice int    `json:"bid_price"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Repository struct {
	DB *gorm.DB
}

func (repo *Repository) AuctionHandler(w http.ResponseWriter, r *http.Request) error {
	adPlacementID := r.URL.Query().Get("ad_placement_id")
	log.Println("Request received for AdPlacementId:", adPlacementID)

	biddingServices := []string{"http://bidding:8081/bid", "http://bidding-2:8081/bid", "http://bidding-3:8081/bid"}

	bidResponses := make(chan BiddingServiceResponse, len(biddingServices))

	var wg sync.WaitGroup

	for _, serviceURL := range biddingServices {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			bidResponse, err := callBiddingService(url)
			if err != nil {
				log.Println("Error calling service:", url, err)
				return
			}
			log.Println("Received bid response from:", url, "AdID:", bidResponse.AdID, "Bid Price:", bidResponse.BidPrice)
			bidResponses <- bidResponse
		}(serviceURL)
	}

	wg.Wait()
	close(bidResponses)

	// Use a boolean to check if we have received any valid bid and avoid using a pointer to the interface
	var bestBid BiddingServiceResponse
	receivedAnyBid := false

	for bid := range bidResponses {
		if !receivedAnyBid || bid.BidPrice > bestBid.BidPrice {
			bestBid = bid
			receivedAnyBid = true
		}
	}

	if !receivedAnyBid {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	bidderModel := &models.Bidder{
		ClientId: &bestBid.AdID,
		BidPrice: &bestBid.BidPrice,
	}

	err := repo.DB.Create(&bidderModel).Error
	if err != nil {
		log.Fatal(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to persist on database."})
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bestBid)

	return nil
}

func callBiddingService(url string) (BiddingServiceResponse, error) {
	client := &http.Client{
		Timeout: time.Millisecond * 200,
	}

	resp, err := client.Get(url)
	if err != nil {
		return BiddingServiceResponse{}, err
	}
	defer resp.Body.Close()

	var bidResponse BiddingServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&bidResponse); err != nil {
		return BiddingServiceResponse{}, err
	}

	return bidResponse, nil
}
