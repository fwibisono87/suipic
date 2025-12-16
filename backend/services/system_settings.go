package services

import (
	"context"
	"strconv"

	"github.com/suipic/backend/repository"
)

type SystemSettingsService struct {
	repo repository.SystemSettingsRepository
}

func NewSystemSettingsService(repo repository.SystemSettingsRepository) *SystemSettingsService {
	return &SystemSettingsService{
		repo: repo,
	}
}

func (s *SystemSettingsService) GetSetting(ctx context.Context, key string) (string, error) {
	return s.repo.Get(ctx, key)
}

func (s *SystemSettingsService) GetAllSettings(ctx context.Context) (map[string]string, error) {
	return s.repo.GetAll(ctx)
}

func (s *SystemSettingsService) UpdateSetting(ctx context.Context, key string, value string) error {
	return s.repo.Set(ctx, key, value)
}

func (s *SystemSettingsService) GetImageProtectionEnabled(ctx context.Context) (bool, error) {
	value, err := s.repo.Get(ctx, "image_protection_enabled")
	if err != nil {
		return false, err
	}

	enabled, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}

	return enabled, nil
}
