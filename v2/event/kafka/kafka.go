package kafka

import(
	"fmt"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"math/rand/v2"
	"github.com/rs/zerolog"

	"github.com/aws/aws-msk-iam-sasl-signer-go/signer"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ProducerWorker struct{
	kafkaConfigurations *KafkaConfigurations
	producer        	*kafka.Producer
	logger 				*zerolog.Logger
}

type ConsumerWorker struct{
	kafkaConfigurations  *KafkaConfigurations
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

func NewProducerWorkerIam(	ctx context.Context,
							region string,
							role string,
							kafkaConfigurations *KafkaConfigurations,
							appLogger *zerolog.Logger) (*ProducerWorker, error) {
	logger := appLogger.With().
						Str("component", "go-core.v2.event.kafka").
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
				Err(err).Msg("erro NewProducer")
		return nil, err
	}

	err = producer.SetOAuthBearerToken(bearerToken)
	if err != nil {
		logger.Error().
				Err(err).Msg("erro SetOAuthBearerToken")
		return nil, err
	}

	return &ProducerWorker{ 
		kafkaConfigurations : kafkaConfigurations,					
		producer : producer,
		logger: &logger,
	}, nil
}

/*
func NewProducerWorker(kafkaConfigurations *KafkaConfigurations
						appLogger *zerolog.Logger) (*ProducerWorker, error) {
	p.logger.Debug().
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
		p.logger.Error().
				Err(err).Send()
		return nil, err
	}
	
	return &ProducerWorker{ kafkaConfigurations : kafkaConfigurations,
							producer : producer,
	}, nil
}
*/

func NewProducerWorkerTX(kafkaConfigurations *KafkaConfigurations,
						appLogger *zerolog.Logger) (*ProducerWorker, error) {
	logger := appLogger.With().
						Str("component", "go-core.v2.event.kafka").
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
	
	return &ProducerWorker{ 
		kafkaConfigurations : kafkaConfigurations,
		producer : producer,
		logger: &logger,
	}, nil
}

// Above Producer a event
func (p *ProducerWorker) Producer(event_topic string, 
									key string,
									kafkaHeader []kafka.Header,
									payload []byte) (error){

	p.logger.Debug().
			Str("func","Producer").Send()

	deliveryChan := make(chan kafka.Event)

	err := p.producer.Produce(&kafka.Message {
												TopicPartition: kafka.TopicPartition{	
												Topic: &event_topic, 
												Partition: kafka.PartitionAny,
											},
												Key:    []byte(key),											
												Value: 	payload, 
												Headers: kafkaHeader,
								},deliveryChan)
	if err != nil {
		p.logger.Error().
				 Err(err).Send()
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		p.logger.Debug().
				Msg("+ ERROR ++ ERROR ++ ERROR +")	
		p.logger.Error().
				Err(m.TopicPartition.Error).Msg("delivery failed")
		p.logger.Debug().
				Msg("+ ERROR ++ ERROR ++  ERROR +")
		
		return m.TopicPartition.Error
	}

	p.logger.Debug().Msg("+ + + + + + + + + + + + + + + + + + + + + + + +")		
	p.logger.Debug().Msg("Delivered message to topic")
	p.logger.Debug().Interface("topic            : ", *m.TopicPartition.Topic).Msg("")
	p.logger.Debug().Interface("key              : ", key ).Msg("")
	p.logger.Printf("kafkaHeader                 : %v\n", kafkaHeader)	
	p.logger.Debug().Interface("partition        : ", m.TopicPartition.Partition).Msg("")
	p.logger.Debug().Interface("offset           : ", m.TopicPartition.Offset).Msg("")
	p.logger.Debug().Msg("+ + + + + + + + + + + + + + + + + + + + + + + +")	

	close(deliveryChan)
	return nil
}

// Above Producer with transaction
func (p *ProducerWorker) InitTransactions(ctx context.Context) error{
	p.logger.Debug().
			Str("func","InitTransactions").Send()

	err := p.producer.InitTransactions(ctx);
	if err != nil {
		p.logger.Error().
				 Err(err).Send()
		return err
	}
	return nil
}

// Above Producer begin transaction
func (p *ProducerWorker) BeginTransaction() error{
	p.logger.Debug().
			Str("func","BeginTransaction").Send()

	err := p.producer.BeginTransaction();
	if err != nil {
		p.logger.Error().
				 Err(err).Send()
		return err
	}
	return nil
}

// Above Producer commit transaction
func (p *ProducerWorker) CommitTransaction(ctx context.Context) error{
	p.logger.Debug().
			Str("func","CommitTransaction").Send()

	err := p.producer.CommitTransaction(ctx);
	if err != nil {
		p.logger.Error().
				 Err(err).Send()
		return err
	}
	return nil
}

// Above Producer abort transaction
func (p *ProducerWorker) AbortTransaction(ctx context.Context) error{
	p.logger.Debug().
			Str("func","AbortTransaction").Send()

	err := p.producer.AbortTransaction(ctx);
	if err != nil {
		p.logger.Error().
				 Err(err).Send()
		return err
	}
	return nil
}

// Above Producer abort transaction
func (p *ProducerWorker) Close(){
	p.logger.Debug().
			Str("func","Close").Send()

	p.producer.Close();
}

func NewConsumerWorker(kafkaConfigurations *KafkaConfigurations,
					   appLogger *zerolog.Logger) (*ConsumerWorker, error) {
	
	logger := appLogger.With().
						Str("component", "go-core.v2.event.kafka").
						Logger()
	logger.Debug().
			Str("func","NewConsumerWorker").Send()

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
		logger.Error().
				Err(err).Send()
		return nil, err
	}

	return &ConsumerWorker{ 
		kafkaConfigurations: kafkaConfigurations,
		consumer: consumer,
		logger: &logger,
	}, nil
}

func (c *ConsumerWorker) Consumer(event_topic []string, messages chan <- Message ) {
	c.logger.Debug().
			Str("func","Consumer").Send()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	defer func() { 
		close(messages)
		c.consumer.Close()
		c.logger.Debug().
				 Msg("Closed consumer waiting please !!!")
	}()

	err := c.consumer.SubscribeTopics(event_topic, nil)
	if err != nil {
		c.logger.Error().
				 Err(err).Send()
	}

	run := true
	for run {
		select {
			case sig := <-sigchan:
				c.logger.Debug().
						Interface("Caught signal terminating: ", sig).Send()
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
					c.logger.Error().
							Interface("kafka.PartitionEOF: ",e).Send()
				case *kafka.Message:
					c.logger.Debug().Msg("+ + + + + + + + + + + + + + + + + + + + + + + +")		

					if e.Headers != nil {
						c.logger.Printf("Headers: %v\n", e.Headers)	
					}
					c.logger.Print("Value : ", string(e.Value))

					headers := extractHeaders(e.Headers)
					msg := Message{
							Header: &headers,
							Payload: string(e.Value),
					}
					messages <- msg

					c.logger.Debug().Msg("+ + + + + + + + + + + + + + + + + + + + + + + +")		
				case kafka.Error:
					c.logger.Error().
							Err(e).
							Msg("kafka.Error")
					if e.Code() == kafka.ErrAllBrokersDown {
						run = false
					}
				default:
					c.logger.Debug().
							Interface("default: ",e).
							Msg("Ignored")
			}
		}
	}
}

// Extract header from kafka.header
func extractHeaders(headers []kafka.Header) map[string]string {
	headerMap := make(map[string]string)
	for _, h := range headers {
		headerMap[string(h.Key)] = string(h.Value)
	}
	return headerMap
}

func (c *ConsumerWorker) Commit(){
	c.logger.Debug().
			Str("func","Commit").Send()
	
	c.consumer.Commit()
}