package producer

import (
	"booking/internal/config"
	"booking/pkg"
	"context"
	"encoding/json"
	"time"
)

type Producer struct {
	producer      producer
	eventUsecases eventUsecases
	errChan       chan error
	log           pkg.Logger
	tickerConfig  int64
}

func NewProducer(producer producer, event eventUsecases, log pkg.Logger, cfg config.InfraConfig) *Producer {
	return &Producer{
		producer:      producer,
		eventUsecases: event,
		log:           log,
		errChan:       make(chan error),
		tickerConfig:  cfg.Ticker,
	}
}

func (p *Producer) GetErrChan() <-chan error {
	return p.errChan
}

func (p *Producer) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(p.tickerConfig) * time.Second)
	defer close(p.errChan)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			events, err := p.eventUsecases.GetEvents(ctx)
			if err != nil {
				p.errChan <- err
				continue
			}
			for _, event := range events {
				kafkaEvent := KafkaEvent{
					EventID:     event.ID,
					EventType:   string(event.EventType),
					OccurredAt:  event.CreatedAt,
					AggregateID: event.AggregateID,
					Payload:     event.Payload,
				}

				payload, err := json.Marshal(kafkaEvent)
				if err != nil {
					p.errChan <- err
					continue
				}

				if err := p.producer.SendEvent(ctx, payload); err != nil {
					p.errChan <- err
					continue
				}

				if err := p.eventUsecases.SetProcessed(ctx, event.ID); err != nil {
					p.errChan <- err
					continue
				}
			}
		}
	}
}

type KafkaEvent struct {
	EventID     string      `json:"event_id"`
	EventType   string      `json:"event_type"`
	OccurredAt  time.Time   `json:"occurred_at"`
	AggregateID string      `json:"aggregate_id"`
	Payload     interface{} `json:"payload"`
}
