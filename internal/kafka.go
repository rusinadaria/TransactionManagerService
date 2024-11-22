package internal

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	// "fmt"
	// "time"
	// "errors"
)

func sendMessage(producer *kafka.Producer, topic string, message string) (interface{}, error) {
	deliveryChan := make(chan kafka.Event)

	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	if err != nil {
		return nil, err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return nil, m.TopicPartition.Error
	}

	return nil, nil
}