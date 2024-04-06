package service

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository"
	"context"
)

type Banners interface {
	AddBanner(ctx context.Context, input BannerAddInput) error
	UpdateBanner(ctx context.Context, input BannerUpdateInput) error
	DeleteBanner(ctx context.Context, bannerId int) error
	GetUserBanner(ctx context.Context, input BannerGetByUserInput) (models.Banner, error)
	GetAllBanners(ctx context.Context) ([]models.AdminBanner, error)
}

type Services struct {
	Banners Banners
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Banners: NewBannersService(repos.Banners),
	}
}
