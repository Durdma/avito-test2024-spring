package service

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository"
	"avito-test2024-spring/pkg/auth"
	"context"
	"net/http"
	"strconv"
	"strings"
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

func (s *UsersService) AddUser(ctx context.Context, input UserAddInput) (string, models.ErrService) {
	user := models.User{
		TagId:   input.TagId,
		IsAdmin: input.IsAdmin,
	}

	if user.TagId < 0 {
		return "", models.NewErrorService(http.StatusBadRequest, "users tag_id must be greater or equal to 0")
	}

	userId, err := s.repo.Create(ctx, user)
	if err != nil {
		return "", models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	accessToken, err := s.tokenManager.NewJWT(strconv.FormatInt(int64(userId), 10), 2*time.Second)
	if err != nil {
		return "", models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	return accessToken, models.ErrService{}
}

func (s *UsersService) UpdateUser(ctx context.Context, input models.User) models.ErrService {
	if input.Id <= 0 {
		return models.NewErrorService(http.StatusBadRequest, "users_id must be greater than 0")
	}

	if input.TagId < 0 {
		return models.NewErrorService(http.StatusBadRequest, "users tag_id must be greater or equal to 0")
	}

	err := s.repo.Update(ctx, input)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return models.NewErrorService(http.StatusNotFound, err.Error())
		}

		return models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	return models.ErrService{}
}

func (s *UsersService) DeleteUser(ctx context.Context, userId int) models.ErrService {
	if userId <= 0 {
		return models.NewErrorService(http.StatusBadRequest, "users_id must be greater than 0")
	}

	err := s.repo.Delete(ctx, userId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return models.NewErrorService(http.StatusNotFound, err.Error())
		}
	}

	return models.ErrService{}
}

func (s *UsersService) GetUserById(ctx context.Context, userId int) (models.User, models.ErrService) {
	if userId <= 0 {
		return models.User{}, models.NewErrorService(http.StatusBadRequest, "users_id must be greater than 0")
	}

	user, err := s.repo.GetUserById(ctx, userId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return models.User{}, models.NewErrorService(http.StatusNotFound, err.Error())
		}

		return models.User{}, models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	return user, models.ErrService{}
}

func (s *UsersService) GetAllUsers(ctx context.Context, tagId int, limit int, offset int) ([]models.User, models.ErrService) {
	if limit < 0 {
		return nil, models.NewErrorService(http.StatusBadRequest, "limit must be greater than 0")
	}

	if offset < 0 {
		return nil, models.NewErrorService(http.StatusBadRequest, "offset must be greater or equal to 0")
	}

	if tagId < 0 {
		return nil, models.NewErrorService(http.StatusBadRequest, "users tag_id must be greater or equal to 0")
	}

	users, err := s.repo.GetAllUsers(ctx, tagId, limit, offset)
	if err != nil {
		return nil, models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	return users, models.ErrService{}
}
