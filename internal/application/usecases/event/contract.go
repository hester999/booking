package event

import (
	"booking/internal/domain"
	"context"
)

type tx interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}

type eventRepo interface {
	GetUnprocessedEvents(ctx context.Context) ([]*domain.Event, error)
	SetProcessedAt(ctx context.Context, id string) error
}
