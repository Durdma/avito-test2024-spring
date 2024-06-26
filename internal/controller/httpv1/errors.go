package httpv1

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Message string `json:"error"`
}

func newErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.AbortWithStatusJSON(statusCode, errorResponse{message})
}
