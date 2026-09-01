package producer

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	gocore_kafka "github.com/eliezerraj/go-core/v3/event/kafka"
)

type Payload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestNewProducerWorker(t *testing.T) {
	dialerConfig := gocore_kafka.DialerConfig{
		Username:   "admin",
		Password:   "admin",
		Protocol:   "SASL_PLAINTEXT",
		Mechanisms: "PLAIN", // <-- For testing use PLAIN.
		Brokers:    "localhost:9092",
	}

	kafkaDialer := gocore_kafka.NewKafkaDialer(dialerConfig)
	if kafkaDialer == nil {
		t.Fatal("KafkaDialer is nil")
	} else {
		t.Logf("KafkaDialer created successfully with brokers: %s", dialerConfig)
	}

	producerConfig := kafkaDialer.ProducerConfig("producer-01")

	producerWorker, err := NewProducerWorker(producerConfig)
	if err != nil {
		t.Fatalf("Failed to create ProducerWorker: %v", err)
	}
	if producerWorker == nil {
		t.Fatal("ProducerWorker is nil")
	}

	key := "producer-01:message-06"
	topic := "topic.go_core_v3.test.producer"
	payload := Payload{ID: 1,Name: "my msg 1 from producer-01"}
	
	payload_bytes, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("failed to marshal payload : %s", err)
	}
	kafkaHeaders := []kafka.Header{}
	kafkaHeaders = append(kafkaHeaders, kafka.Header{Key: "custom-header-id", Value: []byte("MY-CUSTOM-HEADER-001")})
	kafkaHeaders = append(kafkaHeaders, kafka.Header{Key: "tracer-id", Value: []byte("MY-TRACER-TEST-002")})

	err = producerWorker.ProduceMessage(context.Background(), topic, key, kafkaHeaders, payload_bytes)
	if err != nil {
		t.Errorf("Failed to produce message: %v", err)
	} else {
		t.Logf("Message produced successfully to topic %s with key %s", topic, key)
	}

	producerWorker.producer.Flush(5000) // Wait for all messages to be delivered
	producerWorker.producer.Close()

}