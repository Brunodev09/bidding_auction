package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestBiddingHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(BidHandler))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("could not send GET request: %v", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		if resp.ContentLength > 0 {
			t.Fatalf("expected no content, got: %v", resp.Body)
		}
	case http.StatusOK:
		var adObjectResponse AdObject
		if err := json.NewDecoder(resp.Body).Decode(&adObjectResponse); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}
		if adObjectResponse.AdID == uuid.Nil {
			t.Fatalf("expected a valid AdID, got: %v", adObjectResponse.AdID)
		}
		if adObjectResponse.BidPrice < 0 {
			t.Fatalf("expected a positive bid price, got: %v", adObjectResponse.BidPrice)
		}

	default:
		t.Fatalf("unexpected status code: %v", resp.Status)
	}
}
