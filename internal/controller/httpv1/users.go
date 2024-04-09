package httpv1

import (
	"avito-test2024-spring/internal/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	adminsControlUser := api.Group("/users")
	{
		adminsControlUser.POST("/", h.addUser)
		adminsControlUser.GET("/", h.getAllUsers)
		adminsControlUser.DELETE("/:id", h.deleteUser)
	}
}

type UserCreateInput struct {
	IsAdmin bool `json:"is_admin" binding:"required"`
	TagId   int  `json:"tag_id,omitempty"`
}

func (h *Handler) addUser(ctx *gin.Context) {
	var user UserCreateInput
	if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
		h.logger.Error().Err(err).
			Str("method", ctx.Request.Method).
			Str("url", ctx.Request.RequestURI).
			Int("status_code", http.StatusBadRequest).
			Msg("error while unmarshalling json")
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(user)
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

func (h *Handler) getAllUsers(ctx *gin.Context) {}

func (h *Handler) deleteUser(ctx *gin.Context) {}

type refreshInput struct {
	Token string `json:"token" binding:"required"`
}

func (h *Handler) refreshUserSession(ctx *gin.Context) {}
