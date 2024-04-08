package postgresql

import (
	"avito-test2024-spring/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FeaturesRepo struct {
	db *pgxpool.Pool
}

func NewFeaturesRepo(db *pgxpool.Pool) *FeaturesRepo {
	return &FeaturesRepo{
		db: db,
	}
}

func (r *FeaturesRepo) Create(ctx context.Context) error {
	query := `INSERT INTO features default values`

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)
	return err
}

func (r *FeaturesRepo) Delete(ctx context.Context, featureId int) error {
	query := `DELETE FROM features WHERE id=@featureId`
	args := pgx.NamedArgs{
		"featureId": featureId,
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
		return errors.New(fmt.Sprintf("feature with id=%v not found", featureId))
	}

	tx.Commit(ctx)
	return nil
}

func (r *FeaturesRepo) GetAllFeatures(ctx context.Context, limit int, offset int) ([]models.Feature, error) {
	query := `SELECT * FROM features order by id offset @offsetIn limit @limitIn`
	args := pgx.NamedArgs{}
	if limit == 0 {
		query = `SELECT * FROM features order by id offset @offsetIn`
		args = pgx.NamedArgs{
			"offsetIn": offset,
		}
	} else {
		args = pgx.NamedArgs{
			"offsetIn": offset,
			"limitIn":  limit,
		}
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			tx.Commit(ctx)
			return nil, err
		}

		tx.Rollback(ctx)
		return nil, err
	}
	defer rows.Close()

	features := make([]models.Feature, 0)
	for rows.Next() {
		feature := models.Feature{}
		err := rows.Scan(&feature.ID)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		features = append(features, feature)
	}

	tx.Commit(ctx)
	return features, nil
}
