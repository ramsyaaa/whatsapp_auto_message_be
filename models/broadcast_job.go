package models

import (
	"time"
)

type BroadcastJob struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientInfo         string    `gorm:"type:text" json:"client_info"`
	BroadcastPlanAt    time.Time `gorm:"type:timestamp" json:"broadcast_plan_at"`
	BroadcastJobStatus string    `gorm:"type:varchar;default:Waiting" json:"broadcast_job_status"`
	CreatedAt          time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt          time.Time `gorm:"type:timestamp" json:"updated_at"`
	BroadcastCode      string    `gorm:"type:varchar" json:"broadcast_code"`
	BroadcastMessage   string    `gorm:"type:text" json:"broadcast_message"`
}

func (BroadcastJob) TableName() string {
	return "broadcast_jobs"
}
