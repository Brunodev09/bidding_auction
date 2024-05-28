package internal

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type KafkaConsumer struct {
	Id        uuid.UUID
	Topic     string
	Partition int32
	Offset    int64
}

func connectConsumer(brokersUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func EnableQueue() error {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		log.Fatal("KAFKA_BROKER environment variable not set")
		return nil
	}
	consumerConfig := KafkaConsumer{
		Id:        uuid.New(),
		Topic:     "bids",
		Partition: 0,
		Offset:    sarama.OffsetOldest,
	}
	worker, err := connectConsumer([]string{kafkaBroker})
	if err != nil {
		log.Fatal(err)
		return err
	}

	consumer, err := worker.ConsumePartition(consumerConfig.Topic, consumerConfig.Partition, consumerConfig.Offset)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("[%s]-Auction consumer has connected.", consumerConfig.Id)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	msgCount := 0

	doneChannel := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("[%s]-Status: Count %d | Topic %s | Message %s \n",
					consumerConfig.Id, msgCount, string(msg.Topic), string(msg.Value))
			case <-sigchan:
				fmt.Printf("Interruption detected for Auction %s", consumerConfig.Id)
				doneChannel <- struct{}{}
			}
		}
	}()
	<-doneChannel
	fmt.Printf("[%s]-Auction processed %d messages", consumerConfig.Id, msgCount)
	if err := worker.Close(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
