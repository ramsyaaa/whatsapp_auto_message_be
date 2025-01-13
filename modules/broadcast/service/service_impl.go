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

func (s *service) StartBroadcast(ctx context.Context, broadcastID int) (map[string]interface{}, error) {
	return s.repo.StartBroadcast(ctx, broadcastID)
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

func (s *service) UpdateBroadcastStatus(ctx context.Context, broadcastID int, status string) (map[string]interface{}, error) {
	return s.repo.UpdateBroadcastStatus(ctx, broadcastID, status)
}

func (s *service) UpdateRecipientBroadcastStatus(ctx context.Context, broadcastID int, recipientID int, status string) (map[string]interface{}, error) {
	return s.repo.UpdateRecipientBroadcastStatus(ctx, broadcastID, recipientID, status)
}
