package service

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository"
	"avito-test2024-spring/pkg/auth"
	"avito-test2024-spring/pkg/cache"
	"context"
)

type Banners interface {
	AddBanner(ctx context.Context, input BannerAddInput) (int, models.ErrService)
	UpdateBanner(ctx context.Context) models.ErrService
	DeleteBanner(ctx context.Context, bannerId int) models.ErrService
	GetUserBanner(ctx context.Context, featureId int, tagId int, lastRevision bool) (models.Banner, models.ErrService)
	GetAllBanners(ctx context.Context, featureId, tagId, limit, offset int) ([]models.AdminBanner, models.ErrService)
}

type Tags interface {
	AddTag(ctx context.Context) models.ErrService
	DeleteTag(ctx context.Context, tagId int) models.ErrService
	GetAllTags(ctx context.Context, limit int, offset int) ([]models.Tag, models.ErrService)
}

type Features interface {
	AddFeature(ctx context.Context) models.ErrService
	DeleteFeature(ctx context.Context, featureId int) models.ErrService
	GetAllFeatures(ctx context.Context, limit int, offset int) ([]models.Feature, models.ErrService)
}

type Users interface {
	AddUser(ctx context.Context, input UserAddInput) (string, models.ErrService)
	UpdateUser(ctx context.Context, input models.User) models.ErrService
	DeleteUser(ctx context.Context, userId int) models.ErrService
	GetUserById(ctx context.Context, userId int) (models.User, models.ErrService)
	GetAllUsers(ctx context.Context, tagId int, limit int, offset int) ([]models.User, models.ErrService)
}

type Services struct {
	Banners  Banners
	Tags     Tags
	Features Features
	Users    Users
}

func NewServices(repos *repository.Repositories, tokenManager auth.TokenManager, cache cache.Cache) *Services {
	return &Services{
		Banners:  NewBannersService(repos.Banners, cache),
		Tags:     NewTagsService(repos.Tags),
		Features: NewFeaturesService(repos.Features),
		Users:    NewUsersService(repos.Users, tokenManager),
	}
}
