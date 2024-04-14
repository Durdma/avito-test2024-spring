package httpv1

import (
	"avito-test2024-spring/internal/service"
	"avito-test2024-spring/pkg/auth"
	"avito-test2024-spring/pkg/cache"
	"avito-test2024-spring/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	bannersService  service.Banners
	tagsService     service.Tags
	featuresService service.Features
	usersService    service.Users
	logger          *logger.Logs
	tokenManager    auth.TokenManager
	cache           cache.Cache
}

func NewHandler(bannersService service.Banners, tagsService service.Tags,
	featuresService service.Features, usersService service.Users, logger *logger.Logs, tokenManager auth.TokenManager,
	cache cache.Cache) *Handler {
	return &Handler{
		bannersService:  bannersService,
		tagsService:     tagsService,
		featuresService: featuresService,
		usersService:    usersService,
		logger:          logger,
		tokenManager:    tokenManager,
		cache:           cache,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initBannersRoutes(v1)
		h.initTagsFeaturesRoutes(v1)
		h.initUsersRoutes(v1)
	}
}
