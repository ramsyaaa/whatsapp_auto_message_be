package repository

import (
	"context"
)

type BroadcastRepository interface {
	CreateBroadcast(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error)
	ImportRecipient(ctx context.Context, broadcastID int, data []map[string]interface{}) (map[string]interface{}, error)
	ImportPecatuRecipient(ctx context.Context, broadcastID int, data []map[string]interface{}) (map[string]interface{}, error)
	StartBroadcast(ctx context.Context, broadcastId int) (map[string]interface{}, error)
	BroadcastDetail(ctx context.Context, broadcastId int) (map[string]interface{}, error)
	GetBroadcastMessage(ctx context.Context, broadcastID int) (map[string]interface{}, error)
	GetAllRecipientByBroadcastID(ctx context.Context, broadcastID int) (map[string]interface{}, error)
	UpdateBroadcastStatus(ctx context.Context, broadcastID int, status string) (map[string]interface{}, error)
	UpdateRecipientBroadcastStatus(ctx context.Context, recipientID int, broadcastID int, status string) (map[string]interface{}, error)
}
