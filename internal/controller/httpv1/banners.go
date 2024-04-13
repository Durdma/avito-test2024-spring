package httpv1

import (
	"avito-test2024-spring/internal/controller/metrics"
	"avito-test2024-spring/internal/service"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) initBannersRoutes(api *gin.RouterGroup) {
	banners := api.Group("", h.userIdentity) // add auth middleware for admin
	banners.Use(metrics.PrometheusMiddleware())
	{
		banners.POST("", h.bannersAdd)
		banners.PATCH("/:id", h.bannersUpdate)
		banners.DELETE("/:id", h.bannersDelete)
		banners.GET("", h.bannersGetAll)
		banners.GET("/user_banner", h.getUserBanner)
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

//TODO refactor response json like in api docs

func (h *Handler) bannersAdd(ctx *gin.Context) {
	isAdmin := ctx.Value(userCtx).(bool)
	if !isAdmin {
		h.logger.Error().
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusUnauthorized).
			Msg("not admin")
		newErrorResponse(ctx, http.StatusUnauthorized, "not admin")
		return
	}

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
	isAdmin := ctx.Value(userCtx).(bool)
	if !isAdmin {
		h.logger.Error().
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusUnauthorized).
			Msg("not admin")
		newErrorResponse(ctx, http.StatusUnauthorized, "not admin")
		return
	}

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
	isAdmin := ctx.Value(userCtx).(bool)
	if !isAdmin {
		h.logger.Error().
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusUnauthorized).
			Msg("not admin")
		newErrorResponse(ctx, http.StatusUnauthorized, "not admin")
		return
	}

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

func (h *Handler) bannersGetAll(ctx *gin.Context) {
	isAdmin := ctx.Value(userCtx).(bool)
	if !isAdmin {
		h.logger.Error().
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusUnauthorized).
			Msg("not admin")
		newErrorResponse(ctx, http.StatusUnauthorized, "not admin")
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("limit") != "" {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("invalid limit format")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("limit") == "" {
		limit = 0
	}

	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("offset") != "" {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("invalid offset format")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("offset") == "" {
		offset = 0
	}

	tagId, err := strconv.Atoi(ctx.Query("tag_id"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("tag_id") != "" {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("invalid tag_id format")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("tag_id") == "" {
		tagId = 0
	}

	featureId, err := strconv.Atoi(ctx.Query("feature_id"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("feature_id") != "" {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("invalid feature_id format")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("feature_id") == "" {
		featureId = 0
	}

	banners, err := h.bannersService.GetAllBanners(ctx, featureId, tagId, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusInternalServerError).
			Msg("")
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, banners)
}

//var (
//	RequestDuration = prometheus.NewHistogramVec(
//		prometheus.HistogramOpts{
//			Name: "http_request_duration_seconds",
//			Help: "Duration of HTTP requests.",
//		},
//		[]string{"method"},
//	)
//)

func (h *Handler) getUserBanner(ctx *gin.Context) {
	//start := time.Now()
	//defer func() {
	//	method := ctx.Request.Method
	//	elapsed := time.Since(start).Seconds()
	//	RequestDuration.WithLabelValues(method).Observe(elapsed)
	//}()

	tagId, err := strconv.Atoi(ctx.Query("tag_id"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("tag_id") != "" {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("invalid tag_id format")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("tag_id") == "" {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("empty tag_id field")
		newErrorResponse(ctx, http.StatusBadRequest, "empty tag_id field")
		return
	}

	featureId, err := strconv.Atoi(ctx.Query("feature_id"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("feature_id") != "" {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("invalid feature_id format")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("feature_id") == "" {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("empty feature_id")
		newErrorResponse(ctx, http.StatusBadRequest, "empty tag_id field")
		return
	}

	lastRevision := false
	lastRevisionQuery := ctx.Query("use_last_revision")
	if lastRevisionQuery != "" {
		if lastRevisionQuery == "true" {
			lastRevision = true
		} else {
			if lastRevisionQuery != "false" {
				h.logger.Error().Err(err).
					Str("method", ctx.Request.Method).
					Str("url", ctx.Request.RequestURI).
					Int("status_code", http.StatusBadRequest).
					Msg("invalid last_revision format")
				newErrorResponse(ctx, http.StatusBadRequest, "invalid last_revision format")
				return
			}
		}
	}

	banner, err := h.bannersService.GetUserBanner(ctx, featureId, tagId, lastRevision)
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusInternalServerError).
			Msg("")
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, banner)
}
