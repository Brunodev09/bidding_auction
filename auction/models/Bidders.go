package models

import (
	"gorm.io/gorm"
)

type Bidder struct {
	id       uint    `gorm: "primary key;autoIncrement"	json: "id"`
	ClientId *string `json: "client_id"`
	BidPrice *int    `json: "bid_price"`
}

func MigrateEvents(db *gorm.DB) error {
	err := db.AutoMigrate(&Bidder{})
	return err
}
