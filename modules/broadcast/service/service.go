package service

import (
	"context"
)

type BroadcastService interface {
	CreateBroadcast(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error)
	ImportRecipient(ctx context.Context, broadcastID int, data []map[string]interface{}) (map[string]interface{}, error)
	ImportPecatuRecipient(ctx context.Context, broadcastID int, data []map[string]interface{}) (map[string]interface{}, error)
	BroadcastDetail(ctx context.Context, broadcastID int) (map[string]interface{}, error)
	GetBroadcastMessage(ctx context.Context, broadcastID int) (map[string]interface{}, error)
	GetAllRecipientByBroadcastID(ctx context.Context, broadcastID int) (map[string]interface{}, error)
	GetPendingRecipientsByBroadcastID(ctx context.Context, broadcastID int) (map[string]interface{}, error)
	GetPaginatedRecipientsByBroadcastID(ctx context.Context, broadcastID int, page int, limit int, search string) (map[string]interface{}, error)
	UpdateBroadcastStatus(ctx context.Context, broadcastID int, status string) (map[string]interface{}, error)
	UpdateRecipientBroadcastStatus(ctx context.Context, recipientID int, broadcastID int, status string) (map[string]interface{}, error)
	DeleteRecipient(ctx context.Context, recipientID int) (map[string]interface{}, error)
	DeleteBroadcast(ctx context.Context, broadcastID int) (map[string]interface{}, error)
	IsAnyRecipientInBroadcast(ctx context.Context, broadcastID int) (bool, error)
	GetAllBroadcasts(ctx context.Context) ([]map[string]interface{}, error)
}
