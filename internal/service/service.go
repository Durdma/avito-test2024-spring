package service

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository"
	"avito-test2024-spring/pkg/auth"
	"context"
	"time"
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

type Users interface {
	AddUser(ctx context.Context, input UserAddInput) (string, error)
}

type Services struct {
	Banners  Banners
	Tags     Tags
	Features Features
	Users    Users
}

func NewServices(repos *repository.Repositories, tokenManager auth.TokenManager,
	accessTTL time.Duration, refreshTTL time.Duration) *Services {
	return &Services{
		Banners:  NewBannersService(repos.Banners, tokenManager, accessTTL, refreshTTL),
		Tags:     NewTagsService(repos.Tags, tokenManager, accessTTL, refreshTTL),
		Features: NewFeaturesService(repos.Features, tokenManager, accessTTL, refreshTTL),
		Users:    NewUsersService(repos.Users, tokenManager, accessTTL, refreshTTL),
	}
}
