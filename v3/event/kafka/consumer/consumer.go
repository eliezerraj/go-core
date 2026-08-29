package consumer

import(
	"go.uber.org/zap"

	"github.com/eliezerraj/go-core/v3/logger"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ConsumerWorker struct{
	Consumer        	*kafka.Consumer
}

func NewConsumerWorker(kafkaConfig *kafka.ConfigMap) (*ConsumerWorker, error) {
	logger.InfoOutCtx("initializing consumer worker SUCCESSFULLY")

	consumer, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		logger.ErrorOutCtx("failed to create kafka consumer", zap.Error(err))
		return nil, err
	}
	return &ConsumerWorker{
		Consumer: consumer,
	}, nil
}

func (c *ConsumerWorker) SubscribeTopics(topics []string) error {
	logger.InfoOutCtx("subscribing to kafka topics")
	
	err := c.Consumer.SubscribeTopics(topics, nil)
	if err != nil {
		logger.ErrorOutCtx("failed to subscribe to kafka topics", zap.Error(err))
		return err
	}
	return nil
}

// Close closes the consumer and releases any resources associated with it.
func (c *ConsumerWorker) Close() error {
	logger.InfoOutCtx("closing kafka consumer")
	err := c.Consumer.Close()
	if err != nil {
		logger.ErrorOutCtx("failed to close kafka consumer", zap.Error(err))
		return err
	}
	return nil
}

// CommitMessage commits the offset of the provided message to Kafka.
func (c *ConsumerWorker) CommitMessage(msg *kafka.Message) error {
	logger.InfoOutCtx("committing message offset to kafka")
	_, err := c.Consumer.CommitMessage(msg)
	if err != nil {
		logger.ErrorOutCtx("failed to commit message offset", zap.Error(err))
		return err
	}
	return nil
}