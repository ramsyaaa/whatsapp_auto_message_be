package models

import (
	"time"
)

type Message struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RecipientNumber string    `gorm:"type:varchar(20)" json:"recipient_number"`
	MessageContent  string    `gorm:"type:text" json:"message_content"`
	Status          string    `gorm:"type:varchar(20);default:Sent" json:"status"`
	SentAt          time.Time `gorm:"type:timestamp" json:"sent_at"`
	CreatedAt       time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamp" json:"updated_at"`
}

func (Message) TableName() string {
	return "messages"
}
