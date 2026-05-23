package producer

import (
	"booking/internal/domain"
	"context"
)

type producer interface {
	SendEvent(ctx context.Context, event []byte) error
}
type eventUsecases interface {
	GetEvents(ctx context.Context) ([]*domain.Event, error)
	SetProcessed(ctx context.Context, id string) error
}
