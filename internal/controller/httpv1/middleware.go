package httpv1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userRole"
)

func (h *Handler) userIdentity(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(ctx, http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(ctx, http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(ctx, http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}

	userIdStr, err := h.tokenManager.Parse(headerParts[1])
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}

	user, errResponse := h.usersService.GetUserById(ctx, userId)
	if errResponse.Status != 0 {
		if errResponse.Status == http.StatusNotFound {
			newErrorResponse(ctx, http.StatusUnauthorized, "Пользователь не авторизован")
			return
		}

		newErrorResponse(ctx, errResponse.Status, errResponse.Error)
		return
	}

	ctx.Set(userCtx, user.IsAdmin)
}
