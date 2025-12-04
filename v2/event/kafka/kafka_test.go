package kafka

import (
	"os"
	"fmt"
	"context"
	"testing"
	"encoding/json"

	"github.com/rs/zerolog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Payload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var logger = zerolog.New(os.Stdout).
						With().
						Str("component", "testgocore.kafka").
						Logger()

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

	// Print kafka configurations pretty	
	s, _ := json.MarshalIndent(kafkaConfigurations, "", "\t")
	fmt.Println(string(s))

	producer_01, err := NewProducerWorker(&kafkaConfigurations, 
										  &logger)
	if err != nil {
		t.Errorf("failed to create producer : %s", err)
	}

	key := "producer-01:abc-123"
	event_topic := "EVENT.TEST"
	payload := Payload{ ID: 1,
						Name: "my msg 1 from producer-01",
						}
	payload_bytes, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("failed to marshal payload : %s", err)
	}

	kafkaHeaders := []kafka.Header{}
	kafkaHeaders = append(kafkaHeaders, kafka.Header{Key: "custom-header-id", Value: []byte("MY-CUSTOM-HEADER-001")})
	kafkaHeaders = append(kafkaHeaders, kafka.Header{Key: "tracer-id", Value: []byte("MY-TRACER-TEST-002")})

	err = producer_01.Producer( event_topic, 
								key, 
								kafkaHeaders, 
								payload_bytes)
	if err != nil {
		t.Errorf("failed to connect kafka : %s", err)
	}

	key = "producer-01:abc-456"
	payload = Payload{ID: 2, 
					  Name: "my msg 2 with NO header from producer-01",
					}
	payload_bytes, err = json.Marshal(payload)
	if err != nil {
		t.Errorf("failed to marshal payload : %s", err)
	}

	err = producer_01.Producer(event_topic, 
								key, 
								nil, 
								payload_bytes)
	if err != nil {
		t.Errorf("failed to connect kafka : %s", err)
	}

	producer_01.Close()
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

	// Print kafka configurations pretty	
	s, _ := json.MarshalIndent(kafkaConfigurations, "", "\t")
	fmt.Println(string(s))

	producer_02, err := NewProducerWorkerTX(&kafkaConfigurations, 
											&logger)
	if err != nil {
		t.Errorf("failed to create producer : %s", err)
	}

	key := "producer-02:def-567"
	event_topic := "EVENT.TEST"
	payload := Payload{ ID: 3, 
						Name: "my mesg 3 with TX from producer-02",
						}
	payload_bytes, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("failed to marshal payload : %s", err)
	}

	kafkaHeaders := []kafka.Header{}
	kafkaHeaders = append(kafkaHeaders, kafka.Header{Key: "custom-header-id", Value: []byte("MY-CUSTOM-HEADER-001")})
	kafkaHeaders = append(kafkaHeaders, kafka.Header{Key: "tracer-id", Value: []byte("MY-TRACER-TEST-002")})

	producer_02.InitTransactions(context.Background())
	if err != nil {
		t.Errorf("failed to InitTransactions kafka : %s", err)
	}
	producer_02.BeginTransaction()
	if err != nil {
		t.Errorf("failed to InitTransactions kafka : %s", err)
	}

	err = producer_02.Producer(event_topic, 
								key, 
								kafkaHeaders, 
								payload_bytes)
	if err != nil {
		t.Errorf("failed to connect kafka : %s", err)
	}

	key = "producer-02:def-890"
	payload = Payload{	ID: 4, 
						Name: "my msg 4 with TX from producer-02",
					}

	payload_bytes, err = json.Marshal(payload)
	if err != nil {
		t.Errorf("failed to marshal payload : %s", err)
	}

	err = producer_02.Producer(event_topic, 
								key, 
								kafkaHeaders, 
								payload_bytes)
	if err != nil {
		t.Errorf("failed to connect kafka : %s", err)
	}

	producer_02.CommitTransaction(context.Background())
	if err != nil {
		t.Errorf("failed to CommitTransaction kafka : %s", err)
	}

	producer_02.Close()
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
		Brokers3: "b-3.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9096",			 
		Groupid:"GROUP-CORE-TEST",			 
		Partition: 3,      
		ReplicationFactor: 1,
		RequiredAcks:  1,    
	}

	// Print kafka configurations pretty	
	s, _ := json.MarshalIndent(kafkaConfigurations, "", "\t")
	fmt.Println(string(s))

	consumer_01, err := NewConsumerWorker(&kafkaConfigurations,
										  &logger)
	if err != nil {
		t.Errorf("failed to create consumer : %s", err)
	}

	event_topics := []string{"EVENT.TEST"}
	//messages := make(chan string)
	message := make(chan Message)

	go consumer_01.Consumer(event_topics, message)

	for msg := range message {
		t.Logf("************************CONSUMER MSG ****************************")	
		t.Logf("====>>>>> msg.Header: %v", msg.Header)	
		t.Logf("====>>>>> msg.Payload: %v", msg.Payload)
		consumer_01.Commit()
	}
}