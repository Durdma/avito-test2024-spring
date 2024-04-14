package httpv1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) initTagsFeaturesRoutes(api *gin.RouterGroup) {
	tags := api.Group("/tags", h.userIdentity)
	{
		tags.POST("/", h.addTag)
		tags.DELETE("/:id", h.deleteTag)
		tags.GET("/", h.getAllTags)
	}

	features := api.Group("/features", h.userIdentity)
	{
		features.POST("/", h.addFeature)
		features.DELETE("/:id", h.deleteFeature)
		features.GET("/", h.getAllFeatures)
	}
}

// @Summary Creates a new tag
// @Description Создание нового тэга
// @Tags tag
// @ID create-tag
// @Security Bearer
// @Param Authorization header string true "Bearer token for authentication"
// @Success 201 {string} string "Тэг успешно создан"
// @Failure 401 {string} string "Пользователь не авторизован"
// @Failure 403 {string} string "Пользователь не имеет доступа"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /tags [post]
func (h *Handler) addTag(ctx *gin.Context) {
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

	err := h.tagsService.AddTag(ctx)
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusInternalServerError).Msg("")
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusCreated)
}

// @Summary Deletes a tag
// @Description Удаление тэга
// @Tags tag
// @ID delete-tag
// @Security Bearer
// @Param Authorization header string true "Bearer token for authentication"
// @Param id path integer true "Идентификатор тэга"
// @Success 204 {string} string "Тэг успешно удален"
// @Failure 400 {object} errorResponse "Некорректные данные"
// @Failure 401 {string} string "Пользователь не авторизован"
// @Failure 403 {string} string "Пользователь не имеет доступа"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /tags/{id} [delete]
func (h *Handler) deleteTag(ctx *gin.Context) {
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

	tagId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("invalid id format")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.tagsService.DeleteTag(ctx, tagId)
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

// @Summary Получение всех тэгов
// @Description Получение всех тэгов
// @Tags tag
// @ID get-tags
// @Security Bearer
// @Param Authorization header string true "Bearer token for authentication"
// @Param limit query integer false "Лимит"
// @Param offset query integer false "Оффсет"
// @Success 200 {array} models.Tag "OK"
// @Failure 401 {string} string "Пользователь не авторизован"
// @Failure 403 {string} string "Пользователь не имеет доступа"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /tags [get]
func (h *Handler) getAllTags(ctx *gin.Context) {
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

	tags, err := h.tagsService.GetAllTags(ctx, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusInternalServerError).
			Msg("")
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tags)
}

// @Summary Creates a new feature
// @Description Создание новой фичи
// @Tags feature
// @ID create-feature
// @Security Bearer
// @Param Authorization header string true "Bearer token for authentication"
// @Success 201 {string} string "Фича успешно создана"
// @Failure 401 {string} string "Пользователь не авторизован"
// @Failure 403 {string} string "Пользователь не имеет доступа"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /features [post]
func (h *Handler) addFeature(ctx *gin.Context) {
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

	err := h.featuresService.AddFeature(ctx)
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusInternalServerError).Msg("")
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusCreated)
}

// @Summary Deletes a feature
// @Description Удаление фичи
// @Tags feature
// @ID delete-feature
// @Security Bearer
// @Param Authorization header string true "Bearer token for authentication"
// @Param id path integer true "Идентификатор фичи"
// @Success 204 {string} string "Фича успешно удален"
// @Failure 400 {object} errorResponse "Некорректные данные"
// @Failure 401 {string} string "Пользователь не авторизован"
// @Failure 403 {string} string "Пользователь не имеет доступа"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /features/{id} [delete]
func (h *Handler) deleteFeature(ctx *gin.Context) {
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

	featureId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("invalid id format")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.featuresService.DeleteFeature(ctx, featureId)
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

	ctx.Status(http.StatusNoContent)
}

// @Summary Получение всех фич
// @Description Получение всех фич
// @Tags feature
// @ID get-features
// @Security Bearer
// @Param Authorization header string true "Bearer token for authentication"
// @Param limit query integer false "Лимит"
// @Param offset query integer false "Оффсет"
// @Success 200 {array} models.Feature "OK"
// @Failure 401 {string} string "Пользователь не авторизован"
// @Failure 403 {string} string "Пользователь не имеет доступа"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /features [get]
func (h *Handler) getAllFeatures(ctx *gin.Context) {
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

	tags, err := h.featuresService.GetAllFeatures(ctx, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusInternalServerError).
			Msg("")
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tags)
}
