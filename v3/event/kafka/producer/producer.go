package producer

import(
	"go.uber.org/zap"

	"github.com/eliezerraj/go-core/v3/logger"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ProducerWorker struct{
	producer        	*kafka.Producer
}

func NewProducerWorker(kafkaConfig *kafka.ConfigMap) (*ProducerWorker, error) {
	logger.InfoOutCtx("initializing producer worker SUCCESSFULLY")

	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		logger.ErrorOutCtx("failed to create kafka producer", zap.Error(err))
		return nil, err
	}
	return &ProducerWorker{
		producer: producer,
	}, nil
}

func (p *ProducerWorker) ProduceMessage(topic string, 
										key string,
								  		kafkaHeader []kafka.Header,
								  		message []byte) error {
	logger.InfoOutCtx("producing message to kafka topic")

	deliveryChan := make(chan kafka.Event)

	err := p.producer.Produce(&kafka.Message {	TopicPartition: kafka.TopicPartition{	
												Topic: &topic, 
												Partition: kafka.PartitionAny,
											},
												Key:    []byte(key),											
												Value: 	message, 
												Headers: kafkaHeader,
											},deliveryChan)
	if err != nil {
		logger.ErrorOutCtx("failed to produce message to kafka topic", zap.Error(err))
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		logger.ErrorOutCtx("ERROR * ERROR * ERROR * delivery failed", zap.Error(m.TopicPartition.Error))
		return m.TopicPartition.Error
	}

	logger.InfoOutCtx("message delivered to kafka topic", zap.String("topic", *m.TopicPartition.Topic), zap.Int32("partition", m.TopicPartition.Partition), zap.Int64("offset", int64(m.TopicPartition.Offset)))
	
	close(deliveryChan)
	return nil
}

// Close closes the producer worker and releases any resources associated with it.
func (p *ProducerWorker) Close() {
	logger.InfoOutCtx("closing producer worker SUCCESSFULLY")
	p.producer.Close()
}

// Flush flushes the producer worker, ensuring that all messages are sent to the Kafka broker within the specified timeout (in milliseconds).
func (p *ProducerWorker) Flush(timeoutMs int) int {
	logger.InfoOutCtx("flushing producer worker SUCCESSFULLY")
	return p.producer.Flush(timeoutMs)
}

// Commit returns the underlying Kafka producer instance.
func (p *ProducerWorker) Commit() *kafka.Producer {
	return p.producer
}