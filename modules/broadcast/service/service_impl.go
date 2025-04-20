package service

import (
	"context"

	"go_whatsapp/modules/broadcast/repository"
)

type service struct {
	repo repository.BroadcastRepository
}

func NewBroadcastService(repo repository.BroadcastRepository) BroadcastService {
	return &service{repo: repo}
}

func (s *service) CreateBroadcast(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	return s.repo.CreateBroadcast(ctx, data)
}

func (s *service) ImportRecipient(ctx context.Context, broadcastID int, data []map[string]interface{}) (map[string]interface{}, error) {
	return s.repo.ImportRecipient(ctx, broadcastID, data)
}

func (s *service) ImportPecatuRecipient(ctx context.Context, broadcastID int, data []map[string]interface{}) (map[string]interface{}, error) {
	return s.repo.ImportPecatuRecipient(ctx, broadcastID, data)
}

func (s *service) BroadcastDetail(ctx context.Context, broadcastID int) (map[string]interface{}, error) {
	return s.repo.BroadcastDetail(ctx, broadcastID)
}

func (s *service) GetBroadcastMessage(ctx context.Context, broadcastID int) (map[string]interface{}, error) {
	return s.repo.GetBroadcastMessage(ctx, broadcastID)
}

func (s *service) GetAllRecipientByBroadcastID(ctx context.Context, broadcastID int) (map[string]interface{}, error) {
	return s.repo.GetAllRecipientByBroadcastID(ctx, broadcastID)
}

func (s *service) GetPendingRecipientsByBroadcastID(ctx context.Context, broadcastID int) (map[string]interface{}, error) {
	return s.repo.GetPendingRecipientsByBroadcastID(ctx, broadcastID)
}

func (s *service) GetPaginatedRecipientsByBroadcastID(ctx context.Context, broadcastID int, page int, limit int, search string) (map[string]interface{}, error) {
	return s.repo.GetPaginatedRecipientsByBroadcastID(ctx, broadcastID, page, limit, search)
}

func (s *service) UpdateBroadcastStatus(ctx context.Context, broadcastID int, status string) (map[string]interface{}, error) {
	return s.repo.UpdateBroadcastStatus(ctx, broadcastID, status)
}

func (s *service) UpdateRecipientBroadcastStatus(ctx context.Context, recipientID int, broadcastID int, status string) (map[string]interface{}, error) {
	return s.repo.UpdateRecipientBroadcastStatus(ctx, recipientID, broadcastID, status)
}
func (s *service) DeleteRecipient(ctx context.Context, recipientID int) (map[string]interface{}, error) {
	return s.repo.DeleteRecipient(ctx, recipientID)
}

func (s *service) IsAnyRecipientInBroadcast(ctx context.Context, broadcastID int) (bool, error) {
	return s.repo.IsAnyRecipientInBroadcast(ctx, broadcastID)
}

func (s *service) GetAllBroadcasts(ctx context.Context) ([]map[string]interface{}, error) {
	return s.repo.GetAllBroadcasts(ctx)
}
