package event

import "time"

type Event struct {
	ID          string     `db:"id"`
	AggregateID string     `db:"aggregate_id"`
	EventType   string     `db:"event_type"`
	Payload     []byte     `db:"payload"`
	CreatedAt   time.Time  `db:"created_at"`
	ProcessedAt *time.Time `db:"processed_at"`
}
