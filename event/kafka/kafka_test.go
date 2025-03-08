package kafka

import (
	"context"
	"testing"
	"encoding/json"
)

type Payload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestGoCore_Kafka_Producer(t *testing.T){

	kafkaConfigurations := KafkaConfigurations{
		Username: "admin",
		Password: "admin",
		Protocol: "SASL_SSL", //SASL_PLAINTEXT SASL_SSL
		Mechanisms: "SCRAM-SHA-512", //PLAIN SCRAM-SHA-512
		Clientid: "GO-CORE-TEST",
		Brokers1: "b-1.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9096",
		Brokers2: "b-2.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9096",		 
		Brokers3: "b-3.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9096",		 
		Groupid:"GROUP-CORE-TEST",			 
		Partition: 3,      
		ReplicationFactor: 1,
		RequiredAcks:  1,    
	}

	var producerWorker ProducerWorker
	
	producer_01, err := producerWorker.NewProducerWorker(&kafkaConfigurations)
	if err != nil {
		t.Errorf("failed to open database : %s", err)
	}

	key := "abc-1234"
	event_topic := "EVENT.TEST"
	payload := Payload{ID: 1, Name: "my teste"}

	payload_bytes, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("failed to marshal payload : %s", err)
	}

	err = producer_01.Producer(context.Background(), event_topic, key, payload_bytes)
	if err != nil {
		t.Errorf("failed to connect kafka : %s", err)
	}

	key = "abc-4567"
	payload = Payload{ID: 2, Name: "my teste"}
	payload_bytes, err = json.Marshal(payload)
	if err != nil {
		t.Errorf("failed to marshal payload : %s", err)
	}

	err = producer_01.Producer(context.Background(), event_topic, key, payload_bytes)
	if err != nil {
		t.Errorf("failed to connect kafka : %s", err)
	}
}

func TestGoCore_Kafka_ProducerTX(t *testing.T){

	kafkaConfigurations := KafkaConfigurations{
		Username: "admin",
		Password: "admin",
		Protocol: "SASL_SSL", //SASL_PLAINTEXT SASL_SSL
		Mechanisms: "SCRAM-SHA-512", //PLAIN SCRAM-SHA-512
		Clientid: "GO-CORE-TEST",
		Brokers1: "b-1.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9096",
		Brokers2: "b-2.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9096",		 
		Brokers3: "b-3.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9096",		 
		Groupid:"GROUP-CORE-TEST",			 
		Partition: 3,      
		ReplicationFactor: 1,
		RequiredAcks:  1,    
	}

	var producerWorker ProducerWorker
	
	producer_01, err := producerWorker.NewProducerWorkerTX(&kafkaConfigurations)
	if err != nil {
		t.Errorf("failed to open database : %s", err)
	}

	key := "abc-1234"
	event_topic := "EVENT.TEST"
	payload := Payload{ID: 1, Name: "my teste"}

	payload_bytes, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("failed to marshal payload : %s", err)
	}

	producer_01.InitTransactions(context.Background())
	if err != nil {
		t.Errorf("failed to InitTransactions kafka : %s", err)
	}
	producer_01.BeginTransaction()
	if err != nil {
		t.Errorf("failed to InitTransactions kafka : %s", err)
	}

	err = producer_01.Producer(context.Background(), event_topic,key, payload_bytes)
	if err != nil {
		t.Errorf("failed to connect kafka : %s", err)
	}

	key = "abc-4567"
	payload = Payload{ID: 2, Name: "my teste"}
	payload_bytes, err = json.Marshal(payload)
	if err != nil {
		t.Errorf("failed to marshal payload : %s", err)
	}

	err = producer_01.Producer(context.Background(), event_topic, key, payload_bytes)
	if err != nil {
		t.Errorf("failed to connect kafka : %s", err)
	}

	producer_01.CommitTransaction(context.Background())
	if err != nil {
		t.Errorf("failed to CommitTransaction kafka : %s", err)
	}
}

func TestGoCore_Kafka_Consumer(t *testing.T){

	kafkaConfigurations := KafkaConfigurations{
		Username: "admin",
		Password: "admin",
		Protocol: "SASL_SSL",
		Mechanisms: "SCRAM-SHA-512",
		Clientid: "GO-CORE-TEST",
		Brokers1: "b-1.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9096",
		Brokers2: "b-2.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9096",		 
		Brokers3: "b-2.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9096",		 
		Groupid:"GROUP-CORE-TEST",			 
		Partition: 3,      
		ReplicationFactor: 1,
		RequiredAcks:  1,    
	}

	var consumerWorker ConsumerWorker

	consumer_01, err := consumerWorker.NewConsumerWorker(&kafkaConfigurations)
	if err != nil {
		t.Errorf("failed to connect kafka : %s", err)
	}

	event_topics := []string{"EVENT.TEST"}
	messages := make(chan string)

	go consumer_01.Consumer(event_topics, messages)

	for msg := range messages {
		t.Logf("=====>>>>> Received message: %v", msg)
		//consumer_01.Commit()
	}
}