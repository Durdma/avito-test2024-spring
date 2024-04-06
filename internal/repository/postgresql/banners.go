package postgresql

import (
	"avito-test2024-spring/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BannersRepo struct {
	db *pgxpool.Pool
}

func NewBannersRepo(db *pgxpool.Pool) *BannersRepo {
	return &BannersRepo{
		db: db,
	}
}

func (r *BannersRepo) Create(ctx context.Context, banner models.AdminBanner) error {
	return nil
}

func (r *BannersRepo) Update(ctx context.Context, banner models.AdminBanner) error {
	return nil
}

func (r *BannersRepo) Delete(ctx context.Context, bannerId int) error {
	return nil
}

func (r *BannersRepo) GetUserBanner(ctx context.Context, user models.User,
	tagId int, featureId int, lastRevision bool) (models.Banner, error) {
	return models.Banner{}, nil
}

func (r *BannersRepo) GetAllBanners(ctx context.Context, featureId int,
	tagId int, limit int, offset int) ([]models.AdminBanner, error) {
	return nil, nil
}
