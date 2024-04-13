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
	adminCtx            = "isAdmin"
)

func (h *Handler) userIdentity(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(ctx, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(ctx, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(ctx, http.StatusUnauthorized, "token is empty")
		return
	}

	userIdStr, err := h.tokenManager.Parse(headerParts[1])
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, "invalid user_id")
		return
	}

	user, err := h.usersService.GetUserById(ctx, userId)
	if err != nil {
		return
	}

	ctx.Set(userCtx, user.IsAdmin)
}
