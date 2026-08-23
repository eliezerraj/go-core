package kafka

import (
	"github.com/eliezerraj/go-core/v3/logger"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaDialer struct {
	KafkaConfig *kafka.ConfigMap
}

type DialerConfig struct {
    Username		string 
    Password		string 
    Protocol		string
    Mechanisms		string
    Brokers			string
	ClientID		string
	GroupID			string
}

// Base configuration for Kafka Dialer
func NewKafkaDialer(config DialerConfig) *KafkaDialer {
	logger.InfoOutCtx("initializing kafka dialer SUCCESSFULLY")

	kafkaConfig := &kafka.ConfigMap{"bootstrap.servers":	config.Brokers,
									"security.protocol":	config.Protocol, 
									"sasl.mechanisms":      config.Mechanisms, 
									"sasl.username":        config.Username,
									"sasl.password":        config.Password,              
									}

	return &KafkaDialer{
		KafkaConfig: kafkaConfig,
	}
}

// ProducerConfig returns a Kafka producer configuration map based on the provided client ID.
func (d *KafkaDialer) ProducerConfig(clientID string) *kafka.ConfigMap {
	producerConfig := *d.KafkaConfig

	producerConfig["acks"] = "all"
	producerConfig["enable.idempotence"] = true
	producerConfig["retries"] = 5
	producerConfig["max.in.flight.requests.per.connection"] = 5
	producerConfig["compression.type"] = "snappy"
	producerConfig["message.timeout.ms"] = 20000
	producerConfig["client.id"] = clientID
	producerConfig["go.logs.channel.enable"] = true

	return &producerConfig
}

// ConsumerConfig returns a Kafka consumer configuration map based on the provided group ID and client ID.
func (d *KafkaDialer) ConsumerConfig(groupID, clientID string) *kafka.ConfigMap {
	consumerConfig := *d.KafkaConfig
	
	consumerConfig["group.id"] = groupID
	consumerConfig["client.id"] = clientID
	consumerConfig["auto.offset.reset"] = "earliest"
	consumerConfig["enable.auto.commit"] = false
	consumerConfig["go.logs.channel.enable"] = true
	consumerConfig["enable.idempotence"] = true
	consumerConfig["session.timeout.ms"] = 300000
	
	return &consumerConfig
}