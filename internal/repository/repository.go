package repository

import (
	"avito-test2024-spring/internal/models"
	"context"
)

// TODO Как реализовать добавление тегов и фич в БД отдбельной структурой или передавать внутри
type Banners interface {
	Create(ctx context.Context, banner models.AdminBanner) error
	Update(ctx context.Context, banner models.AdminBanner) error
	Delete(ctx context.Context, bannerId int) error
	GetUserBanner(ctx context.Context, user models.User, tagId int, featureId int) (models.Banner, error)
	GetAllBanners(ctx context.Context, featureId int, tagId int, limit int, offset int) ([]models.AdminBanner, error)
}
