package service

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository"
	"avito-test2024-spring/pkg/auth"
	"avito-test2024-spring/pkg/cache"
	"context"
)

type Banners interface {
	AddBanner(ctx context.Context, input BannerAddInput) (int, error)
	UpdateBanner(ctx context.Context) error
	DeleteBanner(ctx context.Context, bannerId int) error
	GetUserBanner(ctx context.Context, featureId int, tagId int, lastRevision bool) (models.Banner, error)
	GetAllBanners(ctx context.Context, featureId, tagId, limit, offset int) ([]models.AdminBanner, error)
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
	UpdateUser(ctx context.Context, input models.User) error
	DeleteUser(ctx context.Context, userId int) error
	GetUserById(ctx context.Context, userId int) (models.User, error)
	GetAllUsers(ctx context.Context, tagId int, limit int, offset int) ([]models.User, error)
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
