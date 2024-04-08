package service

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository"
	"context"
	"errors"
)

type TagsService struct {
	repo repository.Tags
}

func NewTagsService(repo repository.Tags) *TagsService {
	return &TagsService{
		repo: repo,
	}
}

func (s *TagsService) AddTag(ctx context.Context) error {
	return s.repo.Create(ctx)
}

func (s *TagsService) DeleteTag(ctx context.Context, tagId int) error {
	if tagId <= 0 {
		return errors.New("index of tag must be greater than 0")
	}

	return s.repo.Delete(ctx, tagId)
}

func (s *TagsService) GetAllTags(ctx context.Context, limit int, offset int) ([]models.Tag, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be greater than 0")
	}

	if offset < 0 {
		return nil, errors.New("offset must be greater or equal to 0")
	}

	return s.repo.GetAllTags(ctx, limit, offset)
}
