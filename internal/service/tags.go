package service

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository"
	"context"
	"net/http"
	"strings"
)

type TagsService struct {
	repo repository.Tags
}

func NewTagsService(repo repository.Tags) *TagsService {
	return &TagsService{
		repo: repo,
	}
}

func (s *TagsService) AddTag(ctx context.Context) (int, models.ErrService) {

	tagId, err := s.repo.Create(ctx)
	if err != nil {
		return -1, models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	return tagId, models.ErrService{}
}

func (s *TagsService) DeleteTag(ctx context.Context, tagId int) models.ErrService {
	if tagId <= 0 {
		return models.NewErrorService(http.StatusBadRequest, "tag_id must be greater than 0")
	}

	err := s.repo.Delete(ctx, tagId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return models.NewErrorService(http.StatusNotFound, err.Error())
		}

		return models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	return models.ErrService{}
}

func (s *TagsService) GetAllTags(ctx context.Context, limit int, offset int) ([]models.Tag, models.ErrService) {
	if limit < 0 {
		return nil, models.NewErrorService(http.StatusBadRequest, "limit must be greater than 0")
	}

	if offset < 0 {
		return nil, models.NewErrorService(http.StatusBadRequest, "offset must be greater or equal to 0")
	}

	tags, err := s.repo.GetAllTags(ctx, limit, offset)
	if err != nil {
		return nil, models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	return tags, models.ErrService{}
}
