package service

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository"
	"avito-test2024-spring/pkg/auth"
	"context"
	"errors"
	"strconv"
	"time"
)

type UsersService struct {
	repo         repository.Users
	tokenManager auth.TokenManager
}

func NewUsersService(repo repository.Users, tokenManager auth.TokenManager) *UsersService {
	return &UsersService{
		repo:         repo,
		tokenManager: tokenManager,
	}
}

type UserAddInput struct {
	IsAdmin bool
	TagId   int
}

func (s *UsersService) AddUser(ctx context.Context, input UserAddInput) (string, error) {
	user := models.User{
		TagId:   input.TagId,
		IsAdmin: input.IsAdmin,
	}

	if user.TagId < 0 {
		return "", errors.New("users tag id must be greater or equal to 0")
	}

	userId, err := s.repo.Create(ctx, user)
	if err != nil {
		return "", err
	}

	accessToken, err := s.tokenManager.NewJWT(strconv.FormatInt(int64(userId), 10), 2*time.Second)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *UsersService) UpdateUser(ctx context.Context, input models.User) error {
	if input.Id <= 0 {
		return errors.New("users id must be greater than 0")
	}

	if input.TagId < 0 {
		return errors.New("users tag id must be greater or equal to 0")
	}

	return s.repo.Update(ctx, input)
}

func (s *UsersService) DeleteUser(ctx context.Context, userId int) error {
	if userId <= 0 {
		return errors.New("users id must be greater than 0")
	}

	return s.repo.Delete(ctx, userId)
}

func (s *UsersService) GetUserById(ctx context.Context, userId int) (models.User, error) {
	if userId <= 0 {
		return models.User{}, errors.New("users id must be greater than 0")
	}

	return s.repo.GetUserById(ctx, userId)
}

func (s *UsersService) GetAllUsers(ctx context.Context, tagId int, limit int, offset int) ([]models.User, error) {
	if limit < 0 {
		return nil, errors.New("limit must be greater than 0")
	}

	if offset < 0 {
		return nil, errors.New("offset must be greater or equal to 0")
	}

	if tagId < 0 {
		return nil, errors.New("users tag id must be greater or equal to 0")
	}

	return s.repo.GetAllUsers(ctx, tagId, limit, offset)
}
