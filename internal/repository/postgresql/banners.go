package postgresql

import (
	"avito-test2024-spring/internal/models"
	"context"
	"database/sql"
	"errors"
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
	query := `UPDATE banners SET fk_feature_id = @featureId, title = @titleIn, text = @textIn,
    	url = @urlIn, is_active = @isActive, created_at = @createdAt, updated_at = @updatedAt WHERE id = @bannerId`
	args := pgx.NamedArgs{
		"bannerId":  banner.ID,
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

	res, err := tx.Exec(ctx, query, args)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if res.RowsAffected() == 0 {
		tx.Rollback(ctx)
		return errors.New(fmt.Sprintf("banner with id=%v not found", banner.ID))
	}

	tx.Commit(ctx)
	return nil
}

func (r *BannersRepo) Delete(ctx context.Context, bannerId int) error {
	query := `DELETE FROM banners WHERE id=@bannerId`
	args := pgx.NamedArgs{
		"bannerId": bannerId,
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	res, err := tx.Exec(ctx, query, args)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if res.RowsAffected() == 0 {
		tx.Rollback(ctx)
		return errors.New(fmt.Sprintf("banner with id=%v not found", bannerId))
	}

	tx.Commit(ctx)
	return nil
}

func (r *BannersRepo) GetBannerByID(ctx context.Context, bannerId int) (models.AdminBanner, error) {
	var banner models.AdminBanner

	query := `SELECT id, COALESCE(fk_feature_id::bigint, 0), title, text, url, is_active, created_at, updated_at FROM banners WHERE id=@bannerId`
	args := pgx.NamedArgs{
		"bannerId": bannerId,
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return models.AdminBanner{}, err
	}

	err = tx.QueryRow(ctx, query, args).Scan(&banner.ID, &banner.Feature.ID, &banner.Content.Title, &banner.Content.Text,
		&banner.Content.URL, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			tx.Rollback(ctx)
			return models.AdminBanner{}, errors.New(fmt.Sprintf("banner with id=%v not found", bannerId))
		}
		tx.Rollback(ctx)
		return models.AdminBanner{}, err
	}

	tx.Commit(ctx)
	return banner, nil
}

func (r *BannersRepo) GetUserBanner(ctx context.Context, user models.User,
	tagId int, featureId int, lastRevision bool) (models.Banner, error) {
	return models.Banner{}, nil
}

func (r *BannersRepo) GetAllBanners(ctx context.Context, featureId int,
	tagId int, limit int, offset int) ([]models.AdminBanner, error) {
	query := `SELECT id, COALESCE(fk_feature_id::bigint, 0), title, text, url, is_active, created_at, updated_at` +
		` FROM banners`
	args := pgx.NamedArgs{}

	if featureId != 0 && tagId != 0 {
		query += ` JOIN banners_tags ON banners.id = banners_tags.fk_banner_id` +
			` JOIN tags ON banners_tags.fk_tag_id = tags.id` + ` WHERE tags.id = @tagId AND banners.fk_feature_id = @featureId`

		args["tagId"] = tagId
		args["featureId"] = featureId
	}

	if featureId != 0 && tagId == 0 {
		query += ` WHERE banners.fk_feature_id = @featureId`

		args["featureId"] = featureId
	}

	if featureId == 0 && tagId != 0 {
		query += ` JOIN banners_tags ON banners.id = banners_tags.fk_banner_id` +
			` JOIN tags ON banners_tags.fk_tag_id = tags.id` + ` WHERE tags.id = @tagId`

		args["tagId"] = tagId
	}

	query += ` ORDER BY banners.id`

	if offset != 0 {
		query += ` OFFSET @offsetIn`
		args["offsetIn"] = offset
	}

	if limit != 0 {
		query += ` LIMIT @limitIn`
		args["limitIn"] = limit
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			tx.Commit(ctx)
			return nil, nil
		}

		tx.Rollback(ctx)
		return nil, err
	}
	defer rows.Close()

	banners := make([]models.AdminBanner, 0)
	for rows.Next() {
		banner := models.AdminBanner{}
		err := rows.Scan(&banner.ID, &banner.Feature.ID, &banner.Content.Title, &banner.Content.Text,
			&banner.Content.URL, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		tags, err := r.getBannerTags(ctx, tx, banner.ID)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		banner.Tags = tags

		banners = append(banners, banner)
	}

	tx.Commit(ctx)
	return banners, nil
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

func (r *BannersRepo) getBannerTags(ctx context.Context, tx pgx.Tx, bannerId int) ([]models.Tag, error) {
	query := `SELECT fk_tag_id FROM banners_tags WHERE fk_banner_id = @bannerId`
	args := pgx.NamedArgs{
		"bannerId": bannerId,
	}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}
	defer rows.Close()

	tags := make([]models.Tag, 0)
	for rows.Next() {
		fmt.Println("loop")
		var tag models.Tag
		err := rows.Scan(&tag.ID)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, err
}
