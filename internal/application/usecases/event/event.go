package event

import (
	"booking/internal/domain"
	"context"
)

type EventUsecases struct {
	tx        tx
	eventRepo eventRepo
}

func NewEventUsecases(tx tx, eventRepo eventRepo) *EventUsecases {
	return &EventUsecases{
		tx:        tx,
		eventRepo: eventRepo,
	}
}

func (e *EventUsecases) GetEvents(ctx context.Context) ([]*domain.Event, error) {
	var res []*domain.Event
	err := e.tx.Run(ctx, func(ctx context.Context) error {
		var err error
		res, err = e.eventRepo.GetUnprocessedEvents(ctx)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *EventUsecases) SetProcessed(ctx context.Context, id string) error {
	err := e.tx.Run(ctx, func(ctx context.Context) error {

		err := e.eventRepo.SetProcessedAt(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
