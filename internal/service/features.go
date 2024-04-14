package service

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository"
	"context"
	"net/http"
	"strings"
)

type FeaturesService struct {
	repo repository.Features
}

func NewFeaturesService(repo repository.Features) *FeaturesService {
	return &FeaturesService{
		repo: repo,
	}
}

func (s *FeaturesService) AddFeature(ctx context.Context) models.ErrService {
	err := s.repo.Create(ctx)
	if err != nil {
		return models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	return models.ErrService{}
}

func (s *FeaturesService) DeleteFeature(ctx context.Context, featureId int) models.ErrService {
	if featureId <= 0 {
		return models.NewErrorService(http.StatusBadRequest, "feature_id must be greater than 0")
	}

	err := s.repo.Delete(ctx, featureId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return models.NewErrorService(http.StatusNotFound, err.Error())
		}

		return models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	return models.ErrService{}
}

func (s *FeaturesService) GetAllFeatures(ctx context.Context, limit int, offset int) ([]models.Feature, models.ErrService) {
	if limit < 0 {
		return nil, models.NewErrorService(http.StatusBadRequest, "limit must be greater than 0")
	}

	if offset < 0 {
		return nil, models.NewErrorService(http.StatusBadRequest, "offset must be greater or equal to 0")
	}

	features, err := s.repo.GetAllFeatures(ctx, limit, offset)
	if err != nil {
		return nil, models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	return features, models.ErrService{}
}
