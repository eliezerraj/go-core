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

var kafkaConfigurationsIAM = KafkaConfigurations{
		Clientid: "GO-CORE-TEST",
		Protocol: "SASL_SSL", //SASL_PLAINTEXT SASL_SSL
		Mechanisms: "OAUTHBEARER", //PLAIN SCRAM-SHA-512 OAUTHBEARER
		Brokers1: "b-1.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9098",
		Brokers2: "b-2.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9098",		 
		Brokers3: "b-3.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9098",	 
		Groupid:"GROUP-CORE-TEST",			 
		Partition: 3,      
		ReplicationFactor: 1,
		RequiredAcks:  1,    
}

var kafkaConfigurationsPlain = KafkaConfigurations{
		Protocol: "PLAINTEXT",
		Clientid: "GO-CORE-TEST",
		Brokers1: "b-1.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9092",
		Brokers2: "b-2.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9092",		 
		Brokers3: "b-3.mskarch01.x25pj7.c3.kafka.us-east-2.amazonaws.com:9092",	
		Groupid:"GROUP-CORE-TEST",			 
		Partition: 3,      
		ReplicationFactor: 1,
		RequiredAcks:  1,    
}	

var kafkaConfigurationsSSL = KafkaConfigurations{
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

func TestGoCore_Kafka_ProducerIAM(t *testing.T){
	// Print kafka configurations pretty	
	s, _ := json.MarshalIndent(kafkaConfigurationsIAM, "", "\t")
	fmt.Println(string(s))

	region := "us-east-2"
	role := "arn:aws:iam::908671954593:role/iamMskAccessRole-eliezer"

	producer_01, err := NewProducerWorkerIam(context.Background(),
											region,
											role,
											&kafkaConfigurationsIAM,
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

func TestGoCore_Kafka_ProducerNoAuth(t *testing.T){
	// Print kafka configurations pretty	
	s, _ := json.MarshalIndent(kafkaConfigurationsPlain, "", "\t")
	fmt.Println(string(s))

	producer_01, err := NewProducerWorkerNoAuth(&kafkaConfigurationsPlain, 
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

func TestGoCore_Kafka_Producer(t *testing.T){

	// Print kafka configurations pretty	
	s, _ := json.MarshalIndent(kafkaConfigurationsSSL, "", "\t")
	fmt.Println(string(s))

	producer_01, err := NewProducerWorker(&kafkaConfigurationsSSL, 
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

	// Print kafka configurations pretty	
	s, _ := json.MarshalIndent(kafkaConfigurationsSSL, "", "\t")
	fmt.Println(string(s))

	producer_02, err := NewProducerWorkerTX(&kafkaConfigurationsSSL, 
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

	// Print kafka configurations pretty	
	s, _ := json.MarshalIndent(kafkaConfigurationsSSL, "", "\t")
	fmt.Println(string(s))

	consumer_01, err := NewConsumerWorker(&kafkaConfigurationsSSL,
										  &logger)
	if err != nil {
		t.Errorf("failed to create consumer : %s", err)
	}

	//var EVENT = "EVENT.TEST" 
	var EVENT_CDC = "postgres-cdc.public.inventory"

	event_topics := []string{EVENT_CDC}

	message := make(chan Message)

	go consumer_01.Consumer(event_topics, message)

	for msg := range message {
		t.Logf("************************CONSUMER MSG ****************************")	
		t.Logf("*-*-*-*-*-*-*-*-*-*-*-*- WAITING MSG *-*-*-*-*-*-*-*-*-*-*-*-*-*-")	
		t.Logf("====>>>>> msg.Header: %v", msg.Header)	
		t.Logf("====>>>>> msg.Payload: %v", msg.Payload)
		consumer_01.Commit()
	}
}