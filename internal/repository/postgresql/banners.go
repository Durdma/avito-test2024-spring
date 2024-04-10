package postgresql

import (
	"avito-test2024-spring/internal/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
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
	var id int

	query := `INSERT INTO banners (fk_feature_id, title, text, url, is_active, created_at, updated_at) VALUES (
    @featureId, @titleIn, @textIn, @urlIn, @isActive, @createdAt, @updatedAt) RETURNING id`
	args := pgx.NamedArgs{
		"titleIn":   banner.Content.Title,
		"textIn":    banner.Content.Text,
		"urlIn":     banner.Content.URL,
		"isActive":  banner.IsActive,
		"createdAt": banner.CreatedAt,
		"updatedAt": banner.UpdatedAt,
	}

	if banner.Feature.ID == 0 {
		args["featureId"] = sql.NullInt64{}
	} else {
		args["featureId"] = banner.Feature.ID
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	fmt.Println(banner)

	err = tx.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = r.insertIntoBannersTags(ctx, tx, id, banner.Tags)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)
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

func (r *BannersRepo) insertIntoBannersTags(ctx context.Context, tx pgx.Tx, bannerId int, tagsId []models.Tag) error {
	query := `INSERT INTO banners_tags (fk_banner_id, fk_tag_id) VALUES (@banner, @tag)`
	args := pgx.NamedArgs{
		"banner": bannerId,
	}
	if len(tagsId) > 0 {
		for _, t := range tagsId {
			args["tag"] = t.ID

			_, err := tx.Exec(ctx, query, args)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return nil
}
