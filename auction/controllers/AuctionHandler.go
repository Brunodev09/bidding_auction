package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type BiddingServiceResponse struct {
	AdID     string `json:"ad_id"`
	BidPrice int    `json:"bid_price"`
}

func AuctionHandler(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bestBid)
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
