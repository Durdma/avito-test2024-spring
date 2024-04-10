package httpv1

import (
	"avito-test2024-spring/internal/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initBannersRoutes(api *gin.RouterGroup) {
	banners := api.Group("/") // add auth middleware for admin
	{
		banners.POST("/", h.bannersAdd)
		banners.PATCH("/:id", h.bannersUpdate)
		banners.DELETE("/:id", h.bannersDelete)
		banners.GET("/", h.bannersGetAll)
	}

	userBanner := api.Group("/user_banner") // add auth middleware for user
	{
		userBanner.GET("/", h.getUserBanner)
	}
}

type bannersAddContent struct {
	Title string `json:"title" binding:"required"`
	Text  string `json:"text" binding:"required"`
	URL   string `json:"url" binding:"required"`
}

type bannersAddInput struct {
	Tags     []int             `json:"tags_ids" binding:"required"`
	Feature  int               `json:"feature_id" binding:"required"`
	Content  bannersAddContent `json:"content" binding:"required"`
	IsActive bool              `json:"is_active" binding:"required"`
}

func (h *Handler) bannersAdd(ctx *gin.Context) {
	var banner bannersAddInput
	if err := json.NewDecoder(ctx.Request.Body).Decode(&banner); err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("error while unmarshalling json")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := h.bannersService.AddBanner(ctx, service.BannerAddInput{
		Title:    banner.Content.Title,
		Text:     banner.Content.Text,
		URL:      banner.Content.URL,
		Tags:     banner.Tags,
		Feature:  banner.Feature,
		IsActive: banner.IsActive,
	})
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("error while adding to db")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusCreated)
}

type bannersUpdateContent struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
	URL   string `json:"url,omitempty"`
}

type bannersUpdateInput struct {
	Tags     []int                `json:"tags,omitempty"`
	Feature  int                  `json:"feature,omitempty"`
	Content  bannersUpdateContent `json:"content,omitempty"`
	IsActive bool                 `json:"is_active,omitempty"`
}

func (h *Handler) bannersUpdate(ctx *gin.Context) {

}

func (h *Handler) bannersDelete(ctx *gin.Context) {}

func (h *Handler) bannersGetAll(ctx *gin.Context) {}

func (h *Handler) getUserBanner(ctx *gin.Context) {}
