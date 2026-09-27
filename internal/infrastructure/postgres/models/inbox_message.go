package models

import (
	"time"
)

type InboxMessageModel struct {
	ConsumerName string `gorm:"primaryKey"`
	MessageID    string `gorm:"primaryKey"`
	PayloadHash  string
	ReceivedAt   time.Time
	CompletedAt  *time.Time
}

func (InboxMessageModel) TableName() string {
	return "inbox_messages"
}
