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
)

func (h *Handler) initBannersRoutes(api *gin.RouterGroup) {
	banners := api.Group("/banner", h.userIdentity) // add auth middleware for admin
	banners.Use(metrics.PrometheusMiddleware())
	{
		banners.POST("", h.bannersAdd)
		banners.PATCH("/:id", h.bannersUpdate)
		banners.DELETE("/:id", h.bannersDelete)
		banners.GET("", h.bannersGetAll)
	}

	userBanner := api.Group("", h.userIdentity)
	{
		userBanner.GET("/user_banner", h.getUserBanner)
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

// CreateBanner creates a new banner.
//
// @Summary Creates a new banner.
// @Description This endpoint allows an admin to create a new banner.
// @Tags banner
// @ID create-banner
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token for authentication"
// @Param body body bannersAddInput true "Banner creation request"
// @Success 201 {object} int "Banner created successfully"
// @Failure 400 {object} errorResponse "Invalid data provided"
// @Failure 401 {object} errorResponse "Unauthorized access"
// @Failure 403 {object} errorResponse "Forbidden access"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /banner [post]
func (h *Handler) bannersAdd(ctx *gin.Context) {
	isAdmin := ctx.Value(userCtx).(bool)
	if !isAdmin {
		h.logger.Error(ctx, http.StatusForbidden, "Пользователь не имеет доступа")
		newErrorResponse(ctx, http.StatusForbidden, "Пользователь не имеет доступа")
		return
	}

	var banner bannersAddInput
	if err := json.NewDecoder(ctx.Request.Body).Decode(&banner); err != nil {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	bannerId, err := h.bannersService.AddBanner(ctx, service.BannerAddInput{
		Title:    banner.Content.Title,
		Text:     banner.Content.Text,
		URL:      banner.Content.URL,
		Tags:     banner.Tags,
		Feature:  banner.Feature,
		IsActive: banner.IsActive,
	})
	if err.Status != 0 {
		h.logger.Error(ctx, err.Status, err.Error)
		newErrorResponse(ctx, err.Status, err.Error)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"banner_id": bannerId})
}

// @Summary Обновление содержимого баннера
// @Description Этот эндпоинт предназначен для обновления содержимого баннера по его идентификатору.
// @Tags banner
// @ID update-banner
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token for authentication"
// @Param id path integer true "Идентификатор баннера"
// @Param body body service.bannersUpdateInput true "Запрос на обновление баннера"
// @Success 200 {string} string "OK"
// @Failure 400 {object} errorResponse "Некорректные данные"
// @Failure 401 {object} errorResponse "Пользователь не авторизован"
// @Failure 403 {object} errorResponse "Пользователь не имеет доступа"
// @Failure 404 {object} errorResponse "Баннер не найден"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /banner/{id} [patch]
func (h *Handler) bannersUpdate(ctx *gin.Context) {
	isAdmin := ctx.Value(userCtx).(bool)
	if !isAdmin {
		h.logger.Error(ctx, http.StatusForbidden, "Пользователь не имеет доступа")
		newErrorResponse(ctx, http.StatusForbidden, "Пользователь не имеет доступа")
		return
	}

	serviceCtx := context.WithValue(ctx, "request_body", ctx.Request.Body)
	serviceCtx = context.WithValue(serviceCtx, "banner_id", ctx.Param("id"))

	err := h.bannersService.UpdateBanner(serviceCtx)
	if err.Status != 0 {
		h.logger.Error(ctx, err.Status, err.Error)
		newErrorResponse(ctx, err.Status, err.Error)
		return
	}

	ctx.Status(http.StatusOK)
}

// DELETE /banner/{id}
// Удаление баннера по идентификатору
// @Summary Удаление баннера по идентификатору
// @Description Этот эндпоинт предназначен для удаления баннера по его идентификатору.
// @Tags banner
// @ID delete-banner
// @Security Bearer
// @Param Authorization header string true "Bearer token for authentication"
// @Param id path integer true "Идентификатор баннера"
// @Success 204 {string} string "Баннер успешно удален"
// @Failure 400 {object} errorResponse "Некорректные данные"
// @Failure 401 {object} errorResponse "Пользователь не авторизован"
// @Failure 403 {object} errorResponse "Пользователь не имеет доступа"
// @Failure 404 {object} errorResponse "Баннер для тэга не найден"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /banner/{id} [delete]
func (h *Handler) bannersDelete(ctx *gin.Context) {
	isAdmin := ctx.Value(userCtx).(bool)
	if !isAdmin {
		h.logger.Error(ctx, http.StatusForbidden, "Пользователь не имеет доступа")
		newErrorResponse(ctx, http.StatusForbidden, "Пользователь не имеет доступа")
		return
	}

	bannerId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	errResponse := h.bannersService.DeleteBanner(ctx, bannerId)
	if errResponse.Status != 0 {
		h.logger.Error(ctx, errResponse.Status, errResponse.Error)
		newErrorResponse(ctx, errResponse.Status, errResponse.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary Получение всех баннеров c фильтрацией по фиче и/или тегу
// @Tags banner
// @Description Этот эндпоинт предназначен для получения всех баннеров с возможностью фильтрации по идентификатору фичи и/или тега.
// @ID get-banners
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token for authentication"
// @Param feature_id query integer false "Идентификатор фичи"
// @Param tag_id query integer false "Идентификатор тега"
// @Param limit query integer false "Лимит"
// @Param offset query integer false "Оффсет"
// @Success 200 {array} models.AdminBanner "OK"
// @Failure 401 {object} errorResponse "Пользователь не авторизован"
// @Failure 403 {object} errorResponse "Пользователь не имеет доступа"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /banner [get]
func (h *Handler) bannersGetAll(ctx *gin.Context) {
	isAdmin := ctx.Value(userCtx).(bool)
	if !isAdmin {
		h.logger.Error(ctx, http.StatusForbidden, "Пользователь не имеет доступа")
		newErrorResponse(ctx, http.StatusForbidden, "Пользователь не имеет доступа")
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("limit") != "" {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("limit") == "" {
		limit = 0
	}

	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("offset") != "" {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("offset") == "" {
		offset = 0
	}

	tagId, err := strconv.Atoi(ctx.Query("tag_id"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("tag_id") != "" {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("tag_id") == "" {
		tagId = 0
	}

	featureId, err := strconv.Atoi(ctx.Query("feature_id"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("feature_id") != "" {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("feature_id") == "" {
		featureId = 0
	}

	banners, errResponse := h.bannersService.GetAllBanners(ctx, featureId, tagId, limit, offset)
	if errResponse.Status != 0 {
		h.logger.Error(ctx, errResponse.Status, errResponse.Error)
		newErrorResponse(ctx, errResponse.Status, errResponse.Error)
		return
	}

	ctx.JSON(http.StatusOK, banners)
}

// @Summary Получение баннера для пользователя
// @Tags banner
// @Description This endpoint allows a user to get a banner based on their tag and feature ID.
// @ID get-user-banner
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer token for authentication"
// @Param tag_id query integer true "User tag"
// @Param feature_id query integer true "Feature ID"
// @Param use_last_revision query boolean false "Get the latest information" default(false)
// @Success 200 {object} models.Banner "User banner"
// @Failure 400 {object} errorResponse "Invalid data provided"
// @Failure 401 {object} errorResponse "Unauthorized access"
// @Failure 403 {object} errorResponse "Forbidden access"
// @Failure 404 {object} errorResponse "Banner not found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /user_banner [get]
func (h *Handler) getUserBanner(ctx *gin.Context) {
	tagId, err := strconv.Atoi(ctx.Query("tag_id"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("tag_id") != "" {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("tag_id") == "" {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, "empty tag_id field")
		return
	}

	featureId, err := strconv.Atoi(ctx.Query("feature_id"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("feature_id") != "" {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("feature_id") == "" {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
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
				h.logger.Error(ctx, http.StatusBadRequest, err.Error())
				newErrorResponse(ctx, http.StatusBadRequest, "invalid last_revision format")
				return
			}
		}
	}

	banner, errResponse := h.bannersService.GetUserBanner(ctx, featureId, tagId, lastRevision)
	if errResponse.Status != 0 {
		h.logger.Error(ctx, errResponse.Status, errResponse.Error)
		newErrorResponse(ctx, errResponse.Status, errResponse.Error)
		return
	}

	ctx.JSON(http.StatusOK, banner)
}
