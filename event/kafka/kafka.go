package kafka

import(
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