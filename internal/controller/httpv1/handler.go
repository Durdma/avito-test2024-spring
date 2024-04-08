package httpv1

import (
	"avito-test2024-spring/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	bannersService service.Banners
	tagsService    service.Tags
}

func NewHandler(bannersService service.Banners, tagsService service.Tags) *Handler {
	return &Handler{
		bannersService: bannersService,
		tagsService:    tagsService,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initBannersRoutes(v1)
		h.initTagsFeaturesRoutes(v1)
	}
}
