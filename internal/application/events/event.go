package events

import (
	"time"
	"uuid"
)

type Event struct {
	ID            uuid.UUID
	AggregateID   uuid.UUID
	EventType     string
	CorrelationID uuid.UUID
	CausationID   *uuid.UUID
	Payload       []byte
	OccurredAt    time.Time
}
