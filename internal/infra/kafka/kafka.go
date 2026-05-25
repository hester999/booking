package kafka

import (
	"booking/internal/config"
	"booking/pkg"
	"context"
	"net"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
	logger   pkg.Logger
}

func NewKafkaProducer(producerConfig config.BrokerConfig, logger pkg.Logger) (*Producer, error) {

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true                                             // возвращает подтверждение об удачной отправке, обязателен для sync
	cfg.Producer.Return.Errors = true                                                // возвращает ошибку
	cfg.Producer.RequiredAcks = sarama.WaitForAll                                    // Ожидает ack, от всех реплик.  NoResponse - отправил и забыл.  WaitForLocal - ждет только от лидера кластер
	cfg.Producer.Retry.Backoff = time.Duration(producerConfig.Backoff) * time.Second // задержка по экспоненте
	cfg.Producer.Retry.Max = producerConfig.MaxRetry
	cfg.Net.MaxOpenRequests = producerConfig.MaxConnections
	cfg.Producer.Partitioner = sarama.NewHashPartitioner //определяет в какую партицию попадет сообщение

	//producer, err := sarama.NewSyncProducer(strings.Split(address, ","), cfg)
	brokers := []string{
		net.JoinHostPort(producerConfig.Host, producerConfig.Port),
	}

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return &Producer{
		topic:    producerConfig.Topic,
		producer: producer,
		logger:   logger,
	}, nil

}

func (p *Producer) SendEvent(ctx context.Context, event []byte) error {

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(event),
	}
	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		p.logger.Error("Failed to send event to Kafka: %w", err)
	}
	return nil
}
