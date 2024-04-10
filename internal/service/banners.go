package service

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository"
	"avito-test2024-spring/pkg/auth"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"time"
)

type BannersService struct {
	repo         repository.Banners
	tokenManager auth.TokenManager
}

func NewBannersService(repo repository.Banners, tokenManager auth.TokenManager) *BannersService {
	return &BannersService{
		repo:         repo,
		tokenManager: tokenManager,
	}
}

type BannerAddInput struct {
	Title    string
	Text     string
	URL      string
	Tags     []int
	Feature  int
	IsActive bool
}

// TODO add created_at and update_at fields update
func (s *BannersService) AddBanner(ctx context.Context, input BannerAddInput) error {
	var banner models.AdminBanner

	bannerContent := models.Banner{
		Title: input.Title,
		Text:  input.Text,
		URL:   input.URL,
	}

	err := bannerContent.ValidateBanner()
	if err != nil {
		return err
	}

	banner.Content = bannerContent

	err = banner.ValidateAndSetFeature(input.Feature)
	if err != nil {
		return err
	}

	err = banner.ValidateAndSetTags(input.Tags)
	if err != nil {
		return err
	}

	banner.CreatedAt = time.Now()
	banner.UpdatedAt = time.Now()

	return s.repo.Create(ctx, banner)
}

type bannersUpdateContent struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
	URL   string `json:"url,omitempty"`
}

type bannersUpdateInput struct {
	Tags     []int                `json:"tags_ids,omitempty"`
	Feature  int                  `json:"feature_id,omitempty"`
	Content  bannersUpdateContent `json:"content,omitempty"`
	IsActive bool                 `json:"is_active,omitempty"`
}

func (i *bannersUpdateInput) setTags(tags []models.Tag) {
	for _, t := range tags {
		i.Tags = append(i.Tags, t.ID)
	}
}

func (s *BannersService) UpdateBanner(ctx context.Context) error {
	bannerId, err := strconv.Atoi(ctx.Value("banner_id").(string))
	if err != nil {
		return err
	}

	bannerOld, err := s.repo.GetBannerByID(ctx, bannerId)
	if err != nil {
		return err
	}

	bannerInput := bannersUpdateInput{
		Feature: bannerOld.Feature.ID,
		Content: bannersUpdateContent{
			Title: bannerOld.Content.Title,
			Text:  bannerOld.Content.Text,
			URL:   bannerOld.Content.URL,
		},
		IsActive: bannerOld.IsActive,
	}

	bannerInput.setTags(bannerOld.Tags)

	if err := json.NewDecoder(ctx.Value("request_body").(io.Reader)).Decode(&bannerInput); err != nil {
		return err
	}

	var banner models.AdminBanner

	bannerContent := models.Banner{
		Title: bannerInput.Content.Title,
		Text:  bannerInput.Content.Text,
		URL:   bannerInput.Content.URL,
	}

	err = bannerContent.ValidateBanner()
	if err != nil {
		return err
	}

	banner.Content = bannerContent

	err = banner.ValidateAndSetFeature(bannerInput.Feature)
	if err != nil {
		return err
	}

	err = banner.ValidateAndSetTags(bannerInput.Tags)
	if err != nil {
		return err
	}

	banner.CreatedAt = bannerOld.CreatedAt
	banner.UpdatedAt = time.Now()
	banner.ID = bannerId
	banner.IsActive = bannerInput.IsActive

	err = s.repo.Update(ctx, banner)
	if err != nil {
		return err
	}

	return nil
}

func (s *BannersService) DeleteBanner(ctx context.Context, bannerId int) error {
	if bannerId <= 0 {
		return errors.New("banner id must be greater than 0")
	}

	return s.repo.Delete(ctx, bannerId)
}

type BannerGetByUserInput struct {
}

func (s *BannersService) GetUserBanner(ctx context.Context, input BannerGetByUserInput) (models.Banner, error) {
	return models.Banner{}, nil
}

func (s *BannersService) GetAllBanners(ctx context.Context) ([]models.AdminBanner, error) {
	return nil, nil
}
