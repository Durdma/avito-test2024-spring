package httpv1

import (
	"avito-test2024-spring/internal/service"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
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

func (h *Handler) bannersUpdate(ctx *gin.Context) {
	serviceCtx := context.WithValue(ctx, "request_body", ctx.Request.Body)
	serviceCtx = context.WithValue(serviceCtx, "banner_id", ctx.Param("id"))

	err := h.bannersService.UpdateBanner(serviceCtx)
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("error while adding to db")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) bannersDelete(ctx *gin.Context) {
	bannerId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("invalid id format")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.bannersService.DeleteBanner(ctx, bannerId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.logger.Error().Err(err).
				Str("method", ctx.Request.Method).
				Str("url", ctx.Request.RequestURI).
				Int("status_code", http.StatusNotFound).
				Msg("")
			newErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}

		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusInternalServerError).
			Msg("")
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) bannersGetAll(ctx *gin.Context) {}

func (h *Handler) getUserBanner(ctx *gin.Context) {}
