package postgresql

import (
	"avito-test2024-spring/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TagsRepo struct {
	db *pgxpool.Pool
}

func NewTagsRepo(db *pgxpool.Pool) *TagsRepo {
	return &TagsRepo{
		db: db,
	}
}

func (r *TagsRepo) Create(ctx context.Context) error {
	query := `INSERT INTO tags default values`

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

func (r *TagsRepo) Delete(ctx context.Context, tagId int) error {
	query := `DELETE FROM tags WHERE id=@tagId`
	args := pgx.NamedArgs{
		"tagId": tagId,
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
		return errors.New(fmt.Sprintf("tag with id=%v not found", tagId))
	}

	tx.Commit(ctx)
	return nil
}

func (r *TagsRepo) GetAllTags(ctx context.Context, limit int, offset int) ([]models.Tag, error) {
	query := `SELECT * FROM tags order by id offset @offsetIn limit @limitIn`
	args := pgx.NamedArgs{}
	if limit == 0 {
		query = `SELECT * FROM tags order by id offset @offsetIn`
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

	tags := make([]models.Tag, 0)
	for rows.Next() {
		tag := models.Tag{}
		err := rows.Scan(&tag.ID)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		tags = append(tags, tag)
	}

	tx.Commit(ctx)
	return tags, nil
}
