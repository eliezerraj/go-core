package kafka

import (
	"testing"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func TestGoCore_Kafka_Producer(t *testing.T){

	kafkaConfigurations := KafkaConfigurations{
		Username: "",
		Password: "",
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

	var event_topic = "EVENT.TEST"
	deliveryChan := make(chan kafka.Event)
	err = producer_01.producer.Produce(&kafka.Message {TopicPartition: kafka.TopicPartition{	
													Topic: &event_topic, 
													Partition: kafka.PartitionAny,
												},
												Key:    []byte("abc-123"),											
												Value: 	[]byte("teste"), 
												Headers:  []kafka.Header{	
																			{
																				Key: "ACCOUNT",
																				Value: []byte("abc-123"), 
																			},
																			{
																				Key: "RequesId",
																				Value: []byte("teste-123"), 
																			},
																		},
								},deliveryChan)
	if err != nil {
		t.Errorf("failed to connect kafka : %s", err)
	}
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		t.Errorf("TopicPartition : %s", m.TopicPartition.Error)	
	} else {
		t.Logf("topic : %s", *m.TopicPartition.Topic)
		//t.Logf("topic : %s", m.TopicPartition.Partition)
		t.Logf("topic : %s", m.TopicPartition.Offset)		
	}
	close(deliveryChan)
}

func TestGoCore_Kafka_Consumer(t *testing.T){

	kafkaConfigurations := KafkaConfigurations{
		Username: "",
		Password: "",
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

	topics := []string{"EVENT.TEST"}
	err = consumer_01.consumer.SubscribeTopics(topics, nil)
	if err != nil {
		t.Errorf("failed to subscribeTopics : %s", err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	run := true
	for run {
		select {
		case sig := <-sigchan:
			t.Logf("Caught signal terminating: : %s", sig)
			run = false
		default:
			ev := consumer_01.consumer.Poll(100)
			if ev == nil {
				continue
			}
		switch e := ev.(type) {
			case kafka.AssignedPartitions:
				consumer_01.consumer.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				consumer_01.consumer.Unassign()	
			case kafka.PartitionEOF:
				t.Logf("kafka.PartitionEOF: : %s", e)
			case *kafka.Message:
				t.Logf("kafka.Message: : %s", string(e.Value))
				run = false
			case kafka.Error:
				t.Errorf("failed to kafka.Error : %s", e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				t.Logf("default : %s", e)
		}
		}
	}
	consumer_01.consumer.Close()
}