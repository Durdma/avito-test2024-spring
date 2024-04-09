package postgresql

import (
	"avito-test2024-spring/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type UsersRepo struct {
	db *pgxpool.Pool
}

func NewUsersRepo(db *pgxpool.Pool) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) Create(ctx context.Context, user models.User, refreshToken string, expiresAt time.Time) (int, error) {
	var id int

	query := `INSERT INTO users (is_admin, fk_tag_id, refresh_token, expires_at) values (@isAdmin, @tagId, @refreshToken, @expiresAt) returning id`
	args := pgx.NamedArgs{}
	if user.TagId == 0 {
		query = `INSERT INTO users (is_admin) values (@isAdmin) returning id`
		args = pgx.NamedArgs{
			"isAdmin":      user.IsAdmin,
			"refreshToken": refreshToken,
			"expiresAt":    expiresAt,
		}
	} else {
		args = pgx.NamedArgs{
			"isAdmin":      user.IsAdmin,
			"tagId":        user.TagId,
			"refreshToken": refreshToken,
			"expiresAt":    expiresAt,
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

func (r *UsersRepo) Update(ctx context.Context, user models.User, tagId int) error {
	return nil
}

func (r *UsersRepo) Delete(ctx context.Context, userId int) error {
	return nil
}

func (r *UsersRepo) GetUserById(ctx context.Context, userId int) (models.User, error) {
	return models.User{}, nil
}

func (r *UsersRepo) GetAllUsers(ctx context.Context, tagId int, limit int, offset int) ([]models.User, error) {
	return nil, nil
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, userId int) (models.User, error) {
	return models.User{}, nil
}

func (r *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (models.User, error) {
	return models.User{}, nil
}

func (r *UsersRepo) SetSession(ctx context.Context, userId int, refreshToken string, refreshTokenTTL time.Duration) error {
	return nil
}
