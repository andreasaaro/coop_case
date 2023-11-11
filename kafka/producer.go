package kafka

import (
	"coop_case/config"
	"github.com/IBM/sarama"
)

// Producer wraps Kafka Async producer
type Producer interface {
	Close() error
	Input() chan<- *sarama.ProducerMessage
	Successes() <-chan *sarama.ProducerMessage
	Errors() <-chan *sarama.ProducerError
}

func NewSaramaProducer(cfg config.KafkaConfig) (Producer, error) {
	config := sarama.NewConfig()
	brokers := []string{cfg.Brokers}
	saramaClient, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}

	return sarama.NewAsyncProducerFromClient(saramaClient)

}
