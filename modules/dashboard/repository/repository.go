package repository

import (
	"context"
	"time"
)

type DashboardRepository interface {
	GetTotalMessages(ctx context.Context) (int64, error)
	GetTotalBroadcasts(ctx context.Context) (int64, error)
	GetTotalRecipients(ctx context.Context) (int64, error)
	GetMessagesByDate(ctx context.Context, days int) (map[string]int64, error)
	GetBroadcastStatusCounts(ctx context.Context) (map[string]int64, error)
	GetRecentBroadcasts(ctx context.Context, limit int) ([]map[string]interface{}, error)
	GetBroadcastSuccessRate(ctx context.Context) (float64, error)
	GetMessagesByHour(ctx context.Context, date time.Time) (map[int]int64, error)
	GetTopRecipients(ctx context.Context, limit int) ([]map[string]interface{}, error)
}
