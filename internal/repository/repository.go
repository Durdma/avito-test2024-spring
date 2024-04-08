package repository

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository/postgresql"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Banners interface {
	Create(ctx context.Context, banner models.AdminBanner) error
	Update(ctx context.Context, banner models.AdminBanner) error
	Delete(ctx context.Context, bannerId int) error
	GetUserBanner(ctx context.Context, user models.User, tagId int, featureId int, lastRevision bool) (models.Banner, error)
	GetAllBanners(ctx context.Context, featureId int, tagId int, limit int, offset int) ([]models.AdminBanner, error)
}

type Tags interface {
	Create(ctx context.Context) error
	//Update(ctx context.Context, tag models.Tag) error //TODO Нужно ли это в рамках задачи
	Delete(ctx context.Context, tagId int) error
	GetAllTags(ctx context.Context, limit int, offset int) ([]models.Tag, error)
}

type Features interface {
	Create(ctx context.Context) error
	//Update(ctx context.Context, feature models.Feature) error //TODO Нужно ли это в рамках задачи
	Delete(ctx context.Context, featureId int) error
	GetAllFeatures(ctx context.Context, limit int, offset int) ([]models.Feature, error)
}

type Users interface {
	Create(ctx context.Context, tagId int) (string, error) // string - accessToken для использования
	Update(ctx context.Context, userId int, tagId int) error
	Delete(ctx context.Context, userId int) error
	GetUserById(ctx context.Context, userId int) (models.User, error)
	GetAllUsers(ctx context.Context, tagId int, limit int, offset int) ([]models.User, error)
}

type Admins interface {
	Create(ctx context.Context) (string, error)           //string - accessToken для использования
	Update(ctx context.Context, admin models.Admin) error //TODO Нужно ли это в рамках задачи
	Delete(ctx context.Context, adminId int) error
	GetAdminById(ctx context.Context, adminId int) (models.Admin, error)
	GetAllAdmins(ctx context.Context) ([]models.Admin, error)
}

type Repositories struct {
	Banners  Banners
	Tags     Tags
	Features Features
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		Banners:  postgresql.NewBannersRepo(db),
		Tags:     postgresql.NewTagsRepo(db),
		Features: postgresql.NewFeaturesRepo(db),
	}
}
