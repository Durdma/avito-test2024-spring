package httpv1

import (
	"avito-test2024-spring/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	bannersService service.Banners
}

func NewHandler(bannersService service.Banners) *Handler {
	return &Handler{
		bannersService: bannersService,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initBannersRoutes(v1)
	}
}
