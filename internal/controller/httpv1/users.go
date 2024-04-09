package httpv1

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/service"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	adminsControlUser := api.Group("/users")
	{
		adminsControlUser.POST("/", h.addUser)
		adminsControlUser.GET("/", h.getAllUsers)
		adminsControlUser.GET("/:id", h.getUserById)
		adminsControlUser.DELETE("/:id", h.deleteUser)
		adminsControlUser.PATCH("/:id", h.updateUser)
	}
}

type UserInput struct {
	IsAdmin bool `json:"is_admin" binding:"required"`
	TagId   int  `json:"tag_id,omitempty"`
}

func (h *Handler) addUser(ctx *gin.Context) {
	var user UserInput
	if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("error while unmarshalling json")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, err := h.usersService.AddUser(ctx, service.UserAddInput{
		TagId:   user.TagId,
		IsAdmin: user.IsAdmin,
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

	ctx.JSON(http.StatusCreated, map[string]string{"access_token": accessToken})
}

func (h *Handler) updateUser(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("invalid id format")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var userInput UserInput
	if err := json.NewDecoder(ctx.Request.Body).Decode(&userInput); err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("error while unmarshalling json")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user := models.User{
		Id:      userId,
		TagId:   userInput.TagId,
		IsAdmin: userInput.IsAdmin,
	}

	err = h.usersService.UpdateUser(ctx, user)
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

func (h *Handler) getUserById(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("invalid id format")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.usersService.GetUserById(ctx, userId)
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

	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) getAllUsers(ctx *gin.Context) {
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

	users, err := h.usersService.GetAllUsers(ctx, tagId, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusInternalServerError).
			Msg("")
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h *Handler) deleteUser(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("invalid id format")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.usersService.DeleteUser(ctx, userId)
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
