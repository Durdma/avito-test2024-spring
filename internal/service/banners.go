package service

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository"
	"context"
)

type BannersService struct {
	repo repository.Banners
}

func NewBannersService(repo repository.Banners) *BannersService {
	return &BannersService{
		repo: repo,
	}
}

type BannerAddInput struct {
	Title    string
	Text     string
	URL      string
	Tags     []int
	Feature  int
	IsActive bool
}

func (s *BannersService) AddBanner(ctx context.Context, input BannerAddInput) error {
	return nil
}

type BannerUpdateInput struct {
}

func (s *BannersService) UpdateBanner(ctx context.Context, input BannerUpdateInput) error {
	return nil
}

func (s *BannersService) DeleteBanner(ctx context.Context, bannerId int) error {
	return nil
}

type BannerGetByUserInput struct {
}

func (s *BannersService) GetUserBanner(ctx context.Context, input BannerGetByUserInput) (models.Banner, error) {
	return models.Banner{}, nil
}

func (s *BannersService) GetAllBanners(ctx context.Context) ([]models.AdminBanner, error) {
	return nil, nil
}