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

type UsersRepo struct {
	db *pgxpool.Pool
}

func NewUsersRepo(db *pgxpool.Pool) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) Create(ctx context.Context, user models.User) (int, error) {
	var id int

	query := `INSERT INTO users (is_admin, fk_tag_id) values (@isAdmin, @tagId) returning id`
	args := pgx.NamedArgs{}
	if user.TagId == 0 {
		query = `INSERT INTO users (is_admin) values (@isAdmin) returning id`
		args = pgx.NamedArgs{
			"isAdmin": user.IsAdmin,
		}
	} else {
		args = pgx.NamedArgs{
			"isAdmin": user.IsAdmin,
			"tagId":   user.TagId,
		}
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return -1, err
	}

	err = tx.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		tx.Rollback(ctx)
		return -1, err
	}

	tx.Commit(ctx)
	return id, nil
}

func (r *UsersRepo) Update(ctx context.Context, user models.User) error {
	query := `UPDATE users SET fk_tag_id = @tagId, is_admin = @isAdmin WHERE id = @userId`
	args := pgx.NamedArgs{}
	if user.TagId == 0 {
		args["tagId"] = sql.NullInt64{}
	} else {
		args["tagId"] = user.TagId
	}

	args["isAdmin"] = user.IsAdmin
	args["userId"] = user.Id

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
		return errors.New(fmt.Sprintf("user with id=%v not found", user.Id))
	}

	tx.Commit(ctx)
	return nil
}

func (r *UsersRepo) Delete(ctx context.Context, userId int) error {
	query := `DELETE FROM users WHERE id=@userId`
	args := pgx.NamedArgs{
		"userId": userId,
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
		return errors.New(fmt.Sprintf("user with id=%v not found", userId))
	}

	tx.Commit(ctx)
	return nil
}

func (r *UsersRepo) GetUserById(ctx context.Context, userId int) (models.User, error) {
	var user models.User

	query := `SELECT id, COALESCE(fk_tag_id::bigint, 0), is_admin FROM users where id=@userId`
	args := pgx.NamedArgs{
		"userId": userId,
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return models.User{}, err
	}

	err = tx.QueryRow(ctx, query, args).Scan(&user.Id, &user.TagId, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			tx.Rollback(ctx)
			return models.User{}, errors.New(fmt.Sprintf("user with id=%v not found", userId))
		}

		tx.Rollback(ctx)
		return models.User{}, err
	}

	tx.Commit(ctx)
	return user, nil
}

func (r *UsersRepo) GetAllUsers(ctx context.Context, tagId int, limit int, offset int) ([]models.User, error) {
	query := `SELECT id, COALESCE(fk_tag_id::bigint, 0), is_admin FROM users`
	args := pgx.NamedArgs{}

	if tagId != 0 {
		query += ` WHERE fk_tag_id=@tagId`
		args["tagId"] = tagId
	}

	query += ` ORDER BY id`

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
			return nil, err
		}

		tx.Rollback(ctx)
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.Id, &user.TagId, &user.IsAdmin)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		users = append(users, user)
	}

	tx.Commit(ctx)
	return users, nil
}
