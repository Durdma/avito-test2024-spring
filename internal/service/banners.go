package service

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/repository"
	"avito-test2024-spring/pkg/cache"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BannersService struct {
	repo  repository.Banners
	cache cache.Cache
}

func NewBannersService(repo repository.Banners, cache cache.Cache) *BannersService {
	return &BannersService{
		repo:  repo,
		cache: cache,
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

func (s *BannersService) AddBanner(ctx context.Context, input BannerAddInput) (int, models.ErrService) {
	var banner models.AdminBanner

	bannerContent := models.Banner{
		Title: input.Title,
		Text:  input.Text,
		URL:   input.URL,
	}

	err := bannerContent.ValidateBanner()
	if err != nil {
		return -1, models.NewErrorService(http.StatusBadRequest, err.Error())
	}

	banner.Content = bannerContent

	err = banner.ValidateAndSetFeature(input.Feature)
	if err != nil {
		return -1, models.NewErrorService(http.StatusBadRequest, err.Error())
	}

	err = banner.ValidateAndSetTags(input.Tags)
	if err != nil {
		return -1, models.NewErrorService(http.StatusBadRequest, err.Error())
	}

	banner.IsActive = input.IsActive

	banner.CreatedAt = time.Now()
	banner.UpdatedAt = time.Now()

	bannerId, err := s.repo.Create(ctx, banner)
	if err != nil {
		return -1, models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	return bannerId, models.ErrService{}
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

func (s *BannersService) UpdateBanner(ctx context.Context) models.ErrService {
	bannerId, err := strconv.Atoi(ctx.Value("banner_id").(string))
	if err != nil {
		return models.NewErrorService(http.StatusBadRequest, err.Error())
	}

	bannerOld, err := s.repo.GetBannerByID(ctx, bannerId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return models.NewErrorService(http.StatusNotFound, err.Error())
		}
		return models.NewErrorService(http.StatusInternalServerError, err.Error())
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

	if err := json.NewDecoder(ctx.Value("request_body").(io.Reader)).Decode(&bannerInput); err != nil {
		return models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	var banner models.AdminBanner

	bannerContent := models.Banner{
		Title: bannerInput.Content.Title,
		Text:  bannerInput.Content.Text,
		URL:   bannerInput.Content.URL,
	}

	err = bannerContent.ValidateBanner()
	if err != nil {
		return models.NewErrorService(http.StatusBadRequest, err.Error())
	}

	banner.Content = bannerContent

	if bannerInput.Feature == 0 {
		banner.Feature.ID = bannerOld.Feature.ID
	} else {
		err = banner.ValidateAndSetFeature(bannerInput.Feature)
		if err != nil {
			return models.NewErrorService(http.StatusBadRequest, err.Error())
		}
	}

	banner.Tags = bannerOld.Tags

	toDel, err := banner.ValidateAndUpdateTags(bannerInput.Tags)
	if err != nil {
		return models.NewErrorService(http.StatusBadRequest, err.Error())
	}

	banner.CreatedAt = bannerOld.CreatedAt
	banner.UpdatedAt = time.Now()
	banner.ID = bannerId
	banner.IsActive = bannerInput.IsActive

	err = s.repo.Update(ctx, banner, toDel)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return models.NewErrorService(http.StatusNotFound, err.Error())
		}
		return models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	if (bannerOld.IsActive != banner.IsActive) && (banner.IsActive == false) {
		err = s.cache.Delete(banner.ID)
		if err != nil {
			return models.NewErrorService(http.StatusInternalServerError, err.Error())
		}
	}

	return models.ErrService{}
}

func (s *BannersService) DeleteBanner(ctx context.Context, bannerId int) models.ErrService {
	if bannerId <= 0 {
		return models.NewErrorService(http.StatusBadRequest, "banner id must be greater than 0")
	}

	err := s.repo.Delete(ctx, bannerId)
	if err != nil {
		return models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	err = s.cache.Delete(bannerId)
	if err != nil {
		return models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	return models.ErrService{}
}

func (s *BannersService) GetUserBanner(ctx context.Context, featureId int, tagId int, lastRevision bool) (models.Banner, models.ErrService) {
	if tagId < 0 {
		return models.Banner{}, models.NewErrorService(http.StatusBadRequest, "tag_id must be greater or equal to 0")
	}

	if featureId < 0 {
		return models.Banner{}, models.NewErrorService(http.StatusBadRequest, "feature_id must be greater or equal to 0")
	}

	if lastRevision {
		banner, bannerId, err := s.repo.GetUserBanner(ctx, featureId, tagId)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				banner, err = s.cache.Get(tagId, featureId)
				if err != nil {
					if strings.Contains(err.Error(), "not found") {
						return models.Banner{}, models.NewErrorService(http.StatusNotFound, err.Error())
					}
					return models.Banner{}, models.NewErrorService(http.StatusInternalServerError, err.Error())
				}
				return banner, models.ErrService{}
			} else {
				return models.Banner{}, models.NewErrorService(http.StatusInternalServerError, err.Error())
			}
		}

		err = s.cache.Set(banner, tagId, featureId, bannerId)
		if err != nil {
			return banner, models.NewErrorService(http.StatusInternalServerError, err.Error())
		}

		return banner, models.ErrService{}
	} else {
		banner, err := s.cache.Get(tagId, featureId)
		if err != nil {
			if err.Error() == "not found" {
				banner, bannerId, err := s.repo.GetUserBanner(ctx, featureId, tagId)
				if err != nil {
					if strings.Contains(err.Error(), "not found") {
						return models.Banner{}, models.NewErrorService(http.StatusNotFound, err.Error())
					}
					return models.Banner{}, models.NewErrorService(http.StatusInternalServerError, err.Error())
				}

				err = s.cache.Set(banner, tagId, featureId, bannerId)
				if err != nil {
					return banner, models.NewErrorService(http.StatusInternalServerError, err.Error())
				}

				return banner, models.ErrService{}
			} else {
				return models.Banner{}, models.NewErrorService(http.StatusInternalServerError, err.Error())
			}
		}

		return banner, models.ErrService{}
	}
}

func (s *BannersService) GetAllBanners(ctx context.Context, featureId, tagId, limit, offset int) ([]models.AdminBanner, models.ErrService) {
	if limit < 0 {
		return nil, models.NewErrorService(http.StatusBadRequest, "limit must be greater than 0")
	}

	if offset < 0 {
		return nil, models.NewErrorService(http.StatusBadRequest, "offset must be greater or equal to 0")
	}

	if tagId < 0 {
		return nil, models.NewErrorService(http.StatusBadRequest, "tag_id must be greater or equal to 0")
	}

	if featureId < 0 {
		return nil, models.NewErrorService(http.StatusBadRequest, "feature_id must be greater or equal to 0")
	}

	banners, err := s.repo.GetAllBanners(ctx, featureId, tagId, limit, offset)
	if err != nil {
		return nil, models.NewErrorService(http.StatusInternalServerError, err.Error())
	}

	return banners, models.ErrService{}
}
