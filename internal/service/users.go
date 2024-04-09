package service

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository"
	"avito-test2024-spring/pkg/auth"
	"context"
	"fmt"
	"strconv"
	"time"
)

type UsersService struct {
	repo         repository.Users
	tokenManager auth.TokenManager

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUsersService(repo repository.Users, tokenManager auth.TokenManager,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *UsersService {
	return &UsersService{
		repo:            repo,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
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

	userId, err := s.repo.Create(ctx, user, "", time.Now())
	if err != nil {
		return "", err
	}

	tmp := strconv.FormatInt(int64(userId), 10)

	fmt.Println(tmp)

	accessToken, err := s.tokenManager.NewJWT(tmp, 2*time.Second)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
