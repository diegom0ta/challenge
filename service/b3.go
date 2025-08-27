package service

import (
	"challenge/models"
	"challenge/repository"
	"database/sql"
	"time"
)

type B3Service struct {
	repo *repository.B3Repository
}

type B3AggregatedData struct {
	Ticker         string  `json:"ticker"`
	MaxRangeValue  float64 `json:"max_range_value"`
	MaxDailyVolume int     `json:"max_daily_volume"`
}

func NewB3Service(db *sql.DB) *B3Service {
	return &B3Service{
		repo: repository.NewB3Repository(db),
	}
}

func (s *B3Service) GetAggregatedData(ticker string, startDate *time.Time) (*B3AggregatedData, error) {
	maxPrice, maxVolume, err := s.repo.GetAggregatedData(ticker, startDate)
	if err != nil {
		return nil, err
	}

	return &B3AggregatedData{
		Ticker:         ticker,
		MaxRangeValue:  maxPrice,
		MaxDailyVolume: maxVolume,
	}, nil
}

func (s *B3Service) GetAll() ([]*models.B3, error) {
	return nil, nil
}

func (s *B3Service) GetByID(id int) (*models.B3, error) {
	return nil, nil
}
