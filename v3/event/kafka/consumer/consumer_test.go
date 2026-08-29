package consumer

import (
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	gocore_kafka "github.com/eliezerraj/go-core/v3/event/kafka"
)

func TestNewConsumerWorker(t *testing.T) {
	dialerConfig := gocore_kafka.DialerConfig{
		Username:   "admin",
		Password:   "admin",
		Protocol:   "SASL_PLAINTEXT",
		Mechanisms: "PLAIN",
		Brokers:    "localhost:9092",
	}

	kafkaDialer := gocore_kafka.NewKafkaDialer(dialerConfig)
	if kafkaDialer == nil {
		t.Fatal("KafkaDialer is nil")
	}

	t.Logf("KafkaDialer created successfully with brokers: %s", dialerConfig.Brokers)

	consumerConfig := kafkaDialer.ConsumerConfig("consumer-group-01", "consumer-01")

	consumerWorker, err := NewConsumerWorker(consumerConfig)
	if err != nil {
		t.Fatalf("Failed to create ConsumerWorker: %v", err)
	}
	if consumerWorker == nil {
		t.Fatal("ConsumerWorker is nil")
	}
	defer func() {
		if err := consumerWorker.Close(); err != nil {
			t.Errorf("Failed to close consumer: %v", err)
		}
	}()

	topics := []string{"topic.go_core_v3.test.producer"}
	if err = consumerWorker.SubscribeTopics(topics); err != nil {
		t.Fatalf("Failed to subscribe to topics: %v", err)
	}

	t.Logf("Subscribed to topics successfully: %v", topics)

	deadline := time.Now().Add(20 * time.Second)
	messagesRead := 0

	for time.Now().Before(deadline) {
		msg, err := consumerWorker.Consumer.ReadMessage(5 * time.Second)
		if err != nil {
			if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() == kafka.ErrTimedOut {
				t.Log("No message received in the current polling window; continuing")
				continue
			}
			t.Errorf("Consumer error: %v", err)
			break
		}

		if msg == nil || msg.Value == nil {
			t.Log("Received nil message, skipping")
			continue
		}

		messagesRead++
		t.Logf("Received message: %s from topic: %s", string(msg.Value), *msg.TopicPartition.Topic)

		if err = consumerWorker.CommitMessage(msg); err != nil {
			t.Errorf("Failed to commit message: %v", err)
			break
		}

		t.Logf("Committed message offset for topic: %s, partition: %d, offset: %d", *msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)

		if messagesRead >= 1 {
			t.Log("Test consumed at least one message; exiting polling loop")
			return
		}
	}

	if messagesRead == 0 {
		t.Log("No messages were received before timeout; this can happen when the topic is empty or the group has already consumed the message")
	}
}
	