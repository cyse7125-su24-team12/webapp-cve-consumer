package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

var (
	kafkaHost     = os.Getenv("KAFKA_HOST")
	kafkaUser     = os.Getenv("KAFKA_USER")
	kafkaPassword = os.Getenv("KAFKA_PASSWORD")
	topicName   = os.Getenv("TOPIC_NAME")
)

func ConsumerKafka() {
    // Kafka broker address
    brokerList := []string{kafkaHost}
	topic := topicName

    // SASL/PLAIN authentication configuration
    config := sarama.NewConfig()
    config.Net.SASL.Enable = true
    config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
    config.Net.SASL.User = kafkaUser
    config.Net.SASL.Password = kafkaPassword
	
	// Enable debug logging
    sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

    config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	// config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // Start consuming from the oldest message

	// Create a new consumer group
	group, err := sarama.NewConsumerGroup(brokerList, "consumer-group", config)
	if err != nil {
		log.Panicf("Error creating consumer group: %v", err)
	}
	defer func() {
		if err := group.Close(); err != nil {
			log.Panicf("Error closing consumer group: %v", err)
		}
	}()

	// Track errors
	go func() {
		for err := range group.Errors() {
			log.Printf("Error: %v", err)
		}
	}()

	// Setup a signal channel to handle SIGINT and SIGTERM for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		cancel()
		group.Close()
	}()

	consumer := &Consumer{}
	for {
		if err := group.Consume(ctx, []string{topic}, consumer); err != nil {
			log.Printf("Error consuming from group: %v", err)
			break
		}
		if ctx.Err() != nil {
			log.Println("Context error:", ctx.Err())
		}
	}
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	db *sql.DB
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Initialization logic here
	log.Println("Consumer setup")

	// Establish the database connection
	log.Println("Connecting to db...")
	db, err := ConnectDB()
	if err != nil {
		log.Printf("Error connecting to the database: %v", err)
		return err
	}
	log.Println("Connected to db...")
	consumer.db = db

	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	// Cleanup logic here
	log.Println("Consumer cleanup")

	// Close the database connection
	if consumer.db != nil {
		consumer.db.Close()
	}

	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	incr := 0
    // Use the established database connection
    db := consumer.db
	if db == nil {
		return fmt.Errorf("database connection is not established")
	}
    for message := range claim.Messages() {
		log.Printf("Message claimed: topic = %s, partition = %d, offset = %d",
        message.Topic, message.Partition, message.Offset)
		
        SaveToDb(message.Value, db)
		session.MarkMessage(message, "") // Mark message as processed
        incr++
        fmt.Println("message counter:",incr)
	}
	return nil
}