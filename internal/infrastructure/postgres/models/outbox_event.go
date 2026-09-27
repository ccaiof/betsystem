package models

import (
	"time"
	"uuid"
)

type OutboxEventModel struct {
	ID            uuid.UUID `gorm:"primaryKey"`
	AggregateID   uuid.UUID
	EventType     string
	CorrelationID *uuid.UUID
	CausationID   *uuid.UUID
	Payload       []byte `gorm:"type:jsonb"`
	OccurredAt    time.Time
	Attempts      int
	NextAttemptAt time.Time
	LockedBy      *string
	LockedAt      *time.Time
	PublishedAt   *time.Time
	Status        string
}

func (OutboxEventModel) TableName() string {
	return "outbox_events"
}
