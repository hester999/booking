package domain

import (
	"booking/utils"
	"time"
)

type EventType string

const (
	BookingCreated   EventType = "booking_created"
	BookingConfirmed EventType = "booking_confirmed"
	BookingCancelled EventType = "booking_cancelled"
	BookingExpired   EventType = "booking_expired"
)

type Event struct {
	ID          string
	AggregateID string
	EventType   EventType
	Payload     []byte
	CreatedAt   time.Time
	ProcessedAt time.Time
}

func NewEvent(AggregateID string, Type EventType, Payload []byte) *Event {
	return &Event{
		ID:          utils.GetUUID(),
		AggregateID: AggregateID,
		EventType:   Type,
		Payload:     Payload,
	}
}

func NewBookingCreatedEvent(aggregateID string, payload []byte) *Event {
	return NewEvent(aggregateID, BookingCreated, payload)
}
func NewBookingConfirmedEvent(aggregateID string, payload []byte) *Event {
	return NewEvent(aggregateID, BookingConfirmed, payload)
}

func NewBookingCancelledEvent(aggregateID string, payload []byte) *Event {
	return NewEvent(aggregateID, BookingCancelled, payload)
}

func NewBookingExpiredEvent(aggregateID string, payload []byte) *Event {
	return NewEvent(aggregateID, BookingExpired, payload)
}
