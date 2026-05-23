package event

import (
	"booking/internal/domain"
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"
)

type EventRepo struct {
	db      *sqlx.DB
	tx      *trmsqlx.CtxGetter
	builder sq.StatementBuilderType
}

func NewEventRepo(db *sqlx.DB, tx *trmsqlx.CtxGetter) *EventRepo {
	return &EventRepo{
		db:      db,
		tx:      tx,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (e *EventRepo) CreateEvent(ctx context.Context, newEvent *domain.Event) error {
	tx := e.tx.DefaultTrOrDB(ctx, e.db)

	query := `insert into outbox (id,aggregate_id,event_type,payload) values ($1, $2, $3,$4)`

	_, err := tx.ExecContext(ctx, query, newEvent.ID, newEvent.AggregateID, string(newEvent.EventType), newEvent.Payload)

	if err != nil {
		return err
	}
	return nil
}

func (e *EventRepo) GetUnprocessedEvents(ctx context.Context) ([]*domain.Event, error) {
	tx := e.tx.DefaultTrOrDB(ctx, e.db)
	var events []Event

	query := `select id,aggregate_id,event_type,payload, created_at, processed_at from outbox where processed_at IS NULL`
	err := tx.SelectContext(ctx, &events, query)

	if err != nil {
		return nil, err
	}

	result := dtoToDomain(events)
	return result, nil
}

func (e *EventRepo) SetProcessedAt(ctx context.Context, id string) error {
	tx := e.tx.DefaultTrOrDB(ctx, e.db)
	query := `update outbox set processed_at = $1 where id = $2`

	_, err := tx.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}
