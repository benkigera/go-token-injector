package services

import (
	"context"
	"mqqt_go/api/updater/models"
	"mqqt_go/api/updater/repositories"
)

type MeterDataService interface {
	GetLatestMeterData(ctx context.Context) (*models.MeterData, error)
}

type meterDataService struct {
	repo repositories.MeterDataRepository
}

func NewMeterDataService(repo repositories.MeterDataRepository) MeterDataService {
	return &meterDataService{repo: repo}
}

func (s *meterDataService) GetLatestMeterData(ctx context.Context) (*models.MeterData, error) {
	// In a real application, you might add business logic here before returning the data
	// For now, it simply calls the repository
	return s.repo.GetLatestMeterData()
}
