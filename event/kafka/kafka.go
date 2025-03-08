package kafka

import(
	"fmt"
	"context"
	"os"
	"os/signal"
	"syscall"
	"math/rand/v2"

	"github.com/rs/zerolog/log"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ProducerWorker struct{
	kafkaConfigurations *KafkaConfigurations
	producer        	*kafka.Producer
}

type ConsumerWorker struct{
	kafkaConfigurations  *KafkaConfigurations
	consumer        	*kafka.Consumer
}

type KafkaConfigurations struct {
    Username		string 
    Password		string 
    Protocol		string
    Mechanisms		string
    Clientid		string 
    Brokers1		string 
    Brokers2		string 
    Brokers3		string 
	Groupid			string 
	Partition       int
    ReplicationFactor int
    RequiredAcks    int
    Lag             int
    LagCommit       int
}

var childLogger = log.With().Str("go-core", "event.kafka").Logger()

func (p *ProducerWorker) NewProducerWorker(kafkaConfigurations *KafkaConfigurations) (*ProducerWorker, error) {
	childLogger.Debug().Msg("NewProducerWorker")

	kafkaBrokerUrls := 	kafkaConfigurations.Brokers1 + "," + kafkaConfigurations.Brokers2 + "," + kafkaConfigurations.Brokers3
	
	config := &kafka.ConfigMap{	"bootstrap.servers":            kafkaBrokerUrls,
								"security.protocol":            kafkaConfigurations.Protocol, //"SASL_SSL",
								"sasl.mechanisms":              kafkaConfigurations.Mechanisms, //"SCRAM-SHA-256",
								"sasl.username":                kafkaConfigurations.Username,
								"sasl.password":                kafkaConfigurations.Password,
								"acks": 						"all", // acks=0  acks=1 acks=all
								"message.timeout.ms":			5000,
								"retries":						5,
								"retry.backoff.ms":				500,
								"enable.idempotence":			true,
								"go.logs.channel.enable": 		true,                    
								}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		childLogger.Error().Err(err).Msg("Failed to create producer:")
		return nil, err
	}
	
	return &ProducerWorker{ kafkaConfigurations : kafkaConfigurations,
							producer : producer,
	}, nil
}

func (p *ProducerWorker) NewProducerWorkerTX(kafkaConfigurations *KafkaConfigurations) (*ProducerWorker, error) {
	childLogger.Debug().Msg("NewProducerWorkerTX")

	kafkaBrokerUrls := 	kafkaConfigurations.Brokers1 + "," + kafkaConfigurations.Brokers2 + "," + kafkaConfigurations.Brokers3
	
	config := &kafka.ConfigMap{	"bootstrap.servers":            kafkaBrokerUrls,
								"security.protocol":            kafkaConfigurations.Protocol, //"SASL_SSL",
								"sasl.mechanisms":              kafkaConfigurations.Mechanisms, //"SCRAM-SHA-256",
								"sasl.username":                kafkaConfigurations.Username,
								"sasl.password":                kafkaConfigurations.Password,
								"acks": 						"all", // acks=0  acks=1 acks=all
								"message.timeout.ms":			5000,
								"retries":						5,
								"retry.backoff.ms":				500,
								"enable.idempotence":			true,
								"go.logs.channel.enable": 		true, 
								"transactional.id":       		fmt.Sprintf("go-core-trx-%v", rand.IntN(1000)),                      
								}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		childLogger.Error().Err(err).Msg("Failed to create producerTX:")
		return nil, err
	}
	
	return &ProducerWorker{ kafkaConfigurations : kafkaConfigurations,
							producer : producer,
	}, nil
}

func (p *ProducerWorker) Producer(ctx context.Context, 
									event_topic string, 
									key string,
									payload []byte) (error){
	childLogger.Debug().Msg("Producer")

	deliveryChan := make(chan kafka.Event)
	err := p.producer.Produce(&kafka.Message {
												TopicPartition: kafka.TopicPartition{	
												Topic: &event_topic, 
												Partition: kafka.PartitionAny,
											},
												Key:    []byte(key),											
												Value: 	payload, 
												Headers:  []kafka.Header{	
																			{
																				Key: "RequesId",
																				Value: []byte(key), 
																			},
																		},
								},deliveryChan)
	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		childLogger.Debug().Msg("+ ERROR + + ERROR + +  ERROR +")	
		childLogger.Error().Err(m.TopicPartition.Error).Msg("delivery failed")
		childLogger.Debug().Msg("+ ERROR + + ERROR + +  ERROR +")
		
		return m.TopicPartition.Error
	}

	childLogger.Debug().Msg("+ + + + + + + + + + + + + + + + + + + + + + + +")		
	childLogger.Debug().Msg("Delivered message to topic")
	childLogger.Debug().Interface("topic    : ",*m.TopicPartition.Topic).Msg("")
	childLogger.Debug().Interface("key      : ",key).Msg("")
	childLogger.Debug().Interface("partition: ", m.TopicPartition.Partition).Msg("")
	childLogger.Debug().Interface("offset   : ",m.TopicPartition.Offset).Msg("")
	childLogger.Debug().Msg("+ + + + + + + + + + + + + + + + + + + + + + + +")	

	close(deliveryChan)
	return nil
}

func (p *ProducerWorker) InitTransactions(ctx context.Context) error{
	childLogger.Debug().Msg("InitTransactions")

	err := p.producer.InitTransactions(ctx);
	if err != nil {
		childLogger.Error().Err(err).Msg("Failed to InitTransactions")
		return err
	}
	return nil
}

func (p *ProducerWorker) BeginTransaction() error{
	childLogger.Debug().Msg("BeginTransaction")

	err := p.producer.BeginTransaction();
	if err != nil {
		childLogger.Error().Err(err).Msg("Failed to BeginTransaction")
		return err
	}
	return nil
}

func (p *ProducerWorker) CommitTransaction(ctx context.Context) error{
	childLogger.Debug().Msg("CommitTransaction")

	err := p.producer.CommitTransaction(ctx);
	if err != nil {
		childLogger.Error().Err(err).Msg("Failed to CommitTransaction")
		return err
	}
	return nil
}

func (p *ProducerWorker) AbortTransaction(ctx context.Context) error{
	childLogger.Debug().Msg("AbortTransaction")

	err := p.producer.AbortTransaction(ctx);
	if err != nil {
		childLogger.Error().Err(err).Msg("Failed to AbortTransaction")
		return err
	}
	return nil
}

func (c *ConsumerWorker) NewConsumerWorker(kafkaConfigurations *KafkaConfigurations) (*ConsumerWorker, error) {
	childLogger.Debug().Msg("NewConsumerWorker")

	kafkaBrokerUrls := 	kafkaConfigurations.Brokers1 + "," + kafkaConfigurations.Brokers2 + "," + kafkaConfigurations.Brokers3
	
	config := &kafka.ConfigMap{	"bootstrap.servers":            kafkaBrokerUrls,
								"security.protocol":            kafkaConfigurations.Protocol, //"SASL_SSL",
								"sasl.mechanisms":              kafkaConfigurations.Mechanisms, //"SCRAM-SHA-256",
								"sasl.username":                kafkaConfigurations.Username,
								"sasl.password":                kafkaConfigurations.Password,
								"group.id":                     kafkaConfigurations.Groupid,
								"enable.auto.commit":           false, //true,
								"broker.address.family": 		"v4",
								"client.id": 					kafkaConfigurations.Clientid,
								"session.timeout.ms":    		6000,
								"enable.idempotence":			true,
								// "auto.offset.reset":     	"latest", 
								"auto.offset.reset":     		"earliest",  
								}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		childLogger.Error().Err(err).Msg("Failed to create consumer")
		return nil, err
	}

	return &ConsumerWorker{ kafkaConfigurations: kafkaConfigurations,
							consumer: 		consumer,
	}, nil
}

func (c *ConsumerWorker) Consumer(event_topic []string,  messages chan<- string ) {
	childLogger.Debug().Msg("Consumer")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	defer func() { 
		close(messages)
		c.consumer.Close()
		childLogger.Debug().Msg("Closed consumer waiting please !!!")
	}()

	err := c.consumer.SubscribeTopics(event_topic, nil)
	if err != nil {
		childLogger.Error().Err(err).Msg("Failed to subscriber topic")
	}

	run := true
	for run {
		select {
			case sig := <-sigchan:
				childLogger.Debug().Interface("Caught signal terminating: ", sig).Msg("")
				run = false
			default:
				ev := c.consumer.Poll(100)
				if ev == nil {
					continue
				}
			switch e := ev.(type) {
				case kafka.AssignedPartitions:
					c.consumer.Assign(e.Partitions)
				case kafka.RevokedPartitions:
					c.consumer.Unassign()	
				case kafka.PartitionEOF:
					childLogger.Error().Interface("kafka.PartitionEOF: ",e).Msg("")
				case *kafka.Message:
					childLogger.Print("----------------------------------")
					if e.Headers != nil {
						childLogger.Printf("Headers: %v\n", e.Headers)
					}
					childLogger.Print("Value : " ,string(e.Value))

					messages <- string(e.Value)
					childLogger.Print("-----------------------------------")
				case kafka.Error:
					childLogger.Error().Err(e).Msg("kafka.Error")
					if e.Code() == kafka.ErrAllBrokersDown {
						run = false
					}
				default:
					childLogger.Debug().Interface("default: ",e).Msg("Ignored")
			}
		}
	}
}

func (c *ConsumerWorker) Commit(){
	childLogger.Debug().Msg("Commit")
	
	c.consumer.Commit()
}