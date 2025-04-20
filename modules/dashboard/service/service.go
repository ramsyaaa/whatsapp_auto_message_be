package service

import (
	"context"
	"time"
)

type DashboardService interface {
	GetDashboardStats(ctx context.Context) (map[string]interface{}, error)
	GetMessageActivity(ctx context.Context, days int) (map[string]interface{}, error)
	GetBroadcastStatus(ctx context.Context) (map[string]interface{}, error)
	GetRecentBroadcasts(ctx context.Context, limit int) ([]map[string]interface{}, error)
	GetHourlyMessageStats(ctx context.Context, date time.Time) (map[string]interface{}, error)
	GetTopRecipients(ctx context.Context, limit int) ([]map[string]interface{}, error)
}
