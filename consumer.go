package main

import (
	"github.com/IBM/sarama"
	"log"
	"database/sql"
)


func consumeMessage(brokerList []string, config *sarama.Config, topic string, db *sql.DB){
	log.Printf("Consuming messages from topic: %s\n", topic)
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Printf("Failed to start consumer: %s", err)
	}
	defer consumer.Close()
	log.Println("\n Consumer started")
	// log.Println(consumer.Partitions("first-topic"))
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Printf("Failed to get partitions: %s", err)
	}

	for _,partition := range partitions{
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			log.Printf("Failed to start consumer for partition %d: %s", partition, err)
		}
		defer pc.Close()
		log.Printf("Consumer started for partition %d\n", partition)
		
		for message := range pc.Messages(){
			log.Printf("Message: %s\n", message.Value)
		}
	}

}
