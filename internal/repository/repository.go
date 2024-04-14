package repository

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository/postgresql"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO refactor query builders like in users repo
type Banners interface {
	Create(ctx context.Context, banner models.AdminBanner) (int, error)
	Update(ctx context.Context, banner models.AdminBanner, toDel []int) error
	Delete(ctx context.Context, bannerId int) error
	GetBannerByID(ctx context.Context, bannerId int) (models.AdminBanner, error)
	GetUserBanner(ctx context.Context, tagId int, featureId int) (models.Banner, int, error)
	GetAllBanners(ctx context.Context, featureId int, tagId int, limit int, offset int) ([]models.AdminBanner, error)
}

type Tags interface {
	Create(ctx context.Context) (int, error)
	Delete(ctx context.Context, tagId int) error
	GetAllTags(ctx context.Context, limit int, offset int) ([]models.Tag, error)
}

type Features interface {
	Create(ctx context.Context) (int, error)
	Delete(ctx context.Context, featureId int) error
	GetAllFeatures(ctx context.Context, limit int, offset int) ([]models.Feature, error)
}

type Users interface {
	Create(ctx context.Context, user models.User) (int, error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, userId int) error
	GetUserById(ctx context.Context, userId int) (models.User, error)
	GetAllUsers(ctx context.Context, tagId int, limit int, offset int) ([]models.User, error)
}

type Repositories struct {
	Banners  Banners
	Tags     Tags
	Features Features
	Users    Users
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		Banners:  postgresql.NewBannersRepo(db),
		Tags:     postgresql.NewTagsRepo(db),
		Features: postgresql.NewFeaturesRepo(db),
		Users:    postgresql.NewUsersRepo(db),
	}
}
