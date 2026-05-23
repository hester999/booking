package event

import "booking/internal/domain"

func dtoToDomain(events []Event) []*domain.Event {
	if len(events) == 0 {
		return nil
	}
	res := make([]*domain.Event, len(events))
	for i, v := range events {
		res[i] = &domain.Event{
			ID:          v.ID,
			AggregateID: v.AggregateID,
			Payload:     v.Payload,
			EventType:   domain.EventType(v.EventType),
			CreatedAt:   v.CreatedAt,
			ProcessedAt: v.ProcessedAt,
		}
	}
	return res
}
