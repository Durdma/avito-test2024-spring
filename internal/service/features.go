package service

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository"
	"context"
	"errors"
)

type FeaturesService struct {
	repo repository.Features
}

func NewFeaturesService(repo repository.Features) *FeaturesService {
	return &FeaturesService{
		repo: repo,
	}
}

func (s *FeaturesService) AddFeature(ctx context.Context) error {
	return s.repo.Create(ctx)
}

func (s *FeaturesService) DeleteFeature(ctx context.Context, featureId int) error {
	if featureId <= 0 {
		return errors.New("index of feature must be greater than 0")
	}

	return s.repo.Delete(ctx, featureId)
}

func (s *FeaturesService) GetAllFeatures(ctx context.Context, limit int, offset int) ([]models.Feature, error) {
	if limit < 0 {
		return nil, errors.New("limit must be greater than 0")
	}

	if offset < 0 {
		return nil, errors.New("offset must be greater or equal to 0")
	}

	return s.repo.GetAllFeatures(ctx, limit, offset)
}
