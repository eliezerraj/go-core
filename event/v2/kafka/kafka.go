package kafka

import(
	"fmt"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"math/rand/v2"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-msk-iam-sasl-signer-go/signer"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ProducerWorker struct{
	kafkaConfigurations *KafkaConfigurations
	producer        	*kafka.Producer
	logger 				*zerolog.Logger
}

type ConsumerWorker struct{
	kafkaConfigurations *KafkaConfigurations
	consumer        	*kafka.Consumer
	logger 				*zerolog.Logger	
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
}

type Message struct {
	Header 	*map[string]string
	Payload string
}

func (p *ProducerWorker) NewProducerWorkerIam(	ctx context.Context,
												region string,
												role string,
												kafkaConfigurations *KafkaConfigurations) (*ProducerWorker, error) {
	logger := appLogger.With().
						Str("component", "go-core.event.v2.kafka").
						Logger()
	logger.Debug().
			Str("func","NewProducerWorkerIam").Send()

	token, tokenExpirationTime, err := signer.GenerateAuthTokenFromRole(ctx, 
																		region, 
																		role, 
																		"go-core-sts-session")
	if err != nil {
		logger.Error().
			Err(err).Send()
		return nil, err
	}

	seconds := tokenExpirationTime / 1000
	nanoseconds := (tokenExpirationTime % 1000) * 1000000

	bearerToken := kafka.OAuthBearerToken{
		TokenValue: token,
		Expiration: time.Unix(seconds, nanoseconds),
	}

	kafkaBrokerUrls := 	kafkaConfigurations.Brokers1 + "," + kafkaConfigurations.Brokers2 + "," + kafkaConfigurations.Brokers3
	
	config := &kafka.ConfigMap{	"bootstrap.servers":            kafkaBrokerUrls,
								"security.protocol":            kafkaConfigurations.Protocol, //"SASL_SSL",
								"sasl.mechanisms":              kafkaConfigurations.Mechanisms, //"SCRAM-SHA-256",
								"acks": 						"all", // acks=0  acks=1 acks=all
								"message.timeout.ms":			5000,
								"retries":						5,
								"retry.backoff.ms":				500,
								"enable.idempotence":			true,
								"go.logs.channel.enable": 		true,                    
								}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		logger.Error().
			Err(err).Send()
		return nil, err
	}

	err = producer.SetOAuthBearerToken(bearerToken)
	if err != nil {
		logger.Error().
			Err(err).Send()
		return nil, err
	}

	return &ProducerWorker{ kafkaConfigurations : kafkaConfigurations,
							producer : producer,
	}, nil
}

func (p *ProducerWorker) NewProducerWorker(kafkaConfigurations *KafkaConfigurations,
											appLogger *zerolog.Logger) (*ProducerWorker, error) {
	logger := appLogger.With().
						Str("component", "go-core.event.v2.kafka").
						Logger()
	logger.Debug().
			Str("func","NewProducerWorker").Send()

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
		logger.Error().
			Err(err).Send()
		return nil, err
	}
	
	return &ProducerWorker{ kafkaConfigurations : kafkaConfigurations,
							producer : producer,
	}, nil
}

func (p *ProducerWorker) NewProducerWorkerTX(kafkaConfigurations *KafkaConfigurations,
											 appLogger *zerolog.Logger) (*ProducerWorker, error) {
	logger := appLogger.With().
						Str("component", "go-core.event.v2.kafka").
						Logger()
	logger.Debug().
			Str("func","NewProducerWorkerTX").Send()

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
		logger.Error().
			Err(err).Send()
		return nil, err
	}
	
	return &ProducerWorker{ kafkaConfigurations : kafkaConfigurations,
							producer : producer,
	}, nil
}

// Above Producer a event
func (p *ProducerWorker) Producer(event_topic string, 
									key string,
									producer_headers *map[string]string,
									payload []byte) (error){
	childLogger.Debug().Str("func","Producer").Send()

	deliveryChan := make(chan kafka.Event)

	var header []kafka.Header

	if producer_headers != nil {
		for key, value := range *producer_headers {
			h := kafka.Header{
				Key: key,
				Value: []byte(value), 
			}
			header = append(header, h)
		}
	}

	err := p.producer.Produce(&kafka.Message {
												TopicPartition: kafka.TopicPartition{	
												Topic: &event_topic, 
												Partition: kafka.PartitionAny,
											},
												Key:    []byte(key),											
												Value: 	payload, 
												Headers: header,
								},deliveryChan)
	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		childLogger.Debug().Msg("+ ERROR ++ ERROR ++ ERROR +")	
		childLogger.Error().Err(m.TopicPartition.Error).Msg("delivery failed")
		childLogger.Debug().Msg("+ ERROR ++ ERROR ++  ERROR +")
		
		return m.TopicPartition.Error
	}

	childLogger.Debug().Msg("+ + + + + + + + + + + + + + + + + + + + + + + +")		
	childLogger.Debug().Msg("Delivered message to topic")
	childLogger.Debug().Interface("topic            : ", *m.TopicPartition.Topic).Msg("")
	childLogger.Debug().Interface("key              : ", key ).Msg("")
	childLogger.Debug().Interface("header           : ", header ).Msg("")
	childLogger.Debug().Interface("partition        : ", m.TopicPartition.Partition).Msg("")
	childLogger.Debug().Interface("offset           : ", m.TopicPartition.Offset).Msg("")
	childLogger.Debug().Msg("+ + + + + + + + + + + + + + + + + + + + + + + +")	

	close(deliveryChan)
	return nil
}

// Above Producer with transaction
func (p *ProducerWorker) InitTransactions(ctx context.Context) error{
	childLogger.Debug().Str("func","InitTransactions").Send()

	err := p.producer.InitTransactions(ctx);
	if err != nil {
		d.logger.Error().
				Err(err).Send()
		return err
	}
	return nil
}

// Above Producer begin transaction
func (p *ProducerWorker) BeginTransaction() error{
	childLogger.Debug().Str("func","BeginTransaction").Send()

	err := p.producer.BeginTransaction();
	if err != nil {
		d.logger.Error().
				Err(err).Send()
		return err
	}
	return nil
}

// Above Producer commit transaction
func (p *ProducerWorker) CommitTransaction(ctx context.Context) error{
	childLogger.Debug().Str("func","CommitTransaction").Send()

	err := p.producer.CommitTransaction(ctx);
	if err != nil {
		d.logger.Error().
				Err(err).Send()
		return err
	}
	return nil
}

// Above Producer abort transaction
func (p *ProducerWorker) AbortTransaction(ctx context.Context) error{
	childLogger.Debug().Str("func","AbortTransaction").Send()

	err := p.producer.AbortTransaction(ctx);
	if err != nil {
		d.logger.Error().
				Err(err).Send()
		return err
	}
	return nil
}

// Above Producer abort transaction
func (p *ProducerWorker) Close(){
	childLogger.Debug().Str("func","Close").Send()

	p.producer.Close();
}

func (c *ConsumerWorker) NewConsumerWorker(kafkaConfigurations *KafkaConfigurations) (*ConsumerWorker, error) {
	childLogger.Debug().Str("func","NewConsumerWorker").Send()

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
								"auto.offset.reset":     		"earliest",  // "earliest" "latest"
								}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		childLogger.Error().Err(err).Send()
		return nil, err
	}

	return &ConsumerWorker{ kafkaConfigurations: kafkaConfigurations,
							consumer: 		consumer,
	}, nil
}

func (c *ConsumerWorker) Consumer(event_topic []string, messages chan <- Message ) {
	childLogger.Debug().Str("func","Consumer").Send()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	defer func() { 
		close(messages)
		c.consumer.Close()
		childLogger.Debug().Msg("Closed consumer waiting please !!!")
	}()

	err := c.consumer.SubscribeTopics(event_topic, nil)
	if err != nil {
		childLogger.Error().Err(err).Send()
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
					childLogger.Print("...................................")

					if e.Headers != nil {
						childLogger.Printf("Headers: %v\n", e.Headers)	
					}
					childLogger.Print("Value : ", string(e.Value))

					headers := extractHeaders(e.Headers)
					msg := Message{
							Header: &headers,
							Payload: string(e.Value),
					}
					messages <- msg

					childLogger.Print("...................................")
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

// Extract header from kafka.header
func extractHeaders(headers []kafka.Header) map[string]string {
	childLogger.Debug().Str("func","extractHeaders").Send()

	childLogger.Debug().Msg("extractHeaders")
	headerMap := make(map[string]string)
	for _, h := range headers {
		headerMap[string(h.Key)] = string(h.Value)
	}
	return headerMap
}

func (c *ConsumerWorker) Commit(){
	childLogger.Debug().Str("func","Commit").Send()
	
	c.consumer.Commit()
}