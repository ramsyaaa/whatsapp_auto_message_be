package models

import (
	"time"
)

type BroadcastRecipient struct {
	ID                        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BroadcastJobID            uint      `gorm:"type:integer" json:"broadcast_job_id"`
	RecipientName             string    `gorm:"type:varchar" json:"recipient_name"`
	RecipientUniqueIdentifier string    `gorm:"type:varchar" json:"recipient_unique_identifier"`
	WhatsappNumber            string    `gorm:"type:varchar" json:"whatsapp_number"`
	BroadcastedAt             time.Time `gorm:"type:timestamp" json:"broadcasted_at"`
	BroadcastStatus           string    `gorm:"type:varchar;default:Pending" json:"broadcast_status"`
	CreatedAt                 time.Time `gorm:"type:timestamp" json:"created_at"`
}

func (BroadcastRecipient) TableName() string {
	return "broadcast_recipients"
}
