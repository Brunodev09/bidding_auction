package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type AdObject struct {
	AdID     uuid.UUID `json:"ad_id"`
	BidPrice int       `json:"bid_price"`
}

func connectProducer(brokersUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func PushEventToQueue(topic string, message []byte) error {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		log.Fatal("KAFKA_BROKER environment variable not set")
	}
	brokersUrl := []string{kafkaBroker}
	producer, err := connectProducer(brokersUrl)
	if err != nil {
		return err
	}
	defer producer.Close()
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Printf("Producer => Log stored in Topic '%s' on partition '%d' and offset '%d'\n", topic, partition, offset)
	return nil
}
