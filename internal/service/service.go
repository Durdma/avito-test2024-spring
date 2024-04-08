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

type Tags interface {
	AddTag(ctx context.Context) error
	DeleteTag(ctx context.Context, tagId int) error
	GetAllTags(ctx context.Context, limit int, offset int) ([]models.Tag, error)
}

type Features interface {
	AddFeature(ctx context.Context) error
	DeleteFeature(ctx context.Context, featureId int) error
	GetAllFeatures(ctx context.Context, limit int, offset int) ([]models.Feature, error)
}

type Services struct {
	Banners  Banners
	Tags     Tags
	Features Features
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Banners:  NewBannersService(repos.Banners),
		Tags:     NewTagsService(repos.Tags),
		Features: NewFeaturesService(repos.Features),
	}
}
