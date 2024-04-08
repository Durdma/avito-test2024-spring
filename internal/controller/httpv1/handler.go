package httpv1

import (
	"avito-test2024-spring/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	bannersService  service.Banners
	tagsService     service.Tags
	featuresService service.Features
	logger          zerolog.Logger
}

func NewHandler(bannersService service.Banners, tagsService service.Tags, featuresService service.Features, logger zerolog.Logger) *Handler {
	return &Handler{
		bannersService:  bannersService,
		tagsService:     tagsService,
		featuresService: featuresService,
		logger:          logger,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initBannersRoutes(v1)
		h.initTagsFeaturesRoutes(v1)
	}
}
