package controller

import (
	docs "avito-test2024-spring/docs"
	"avito-test2024-spring/internal/controller/httpv1"
	"avito-test2024-spring/internal/controller/metrics"
	"avito-test2024-spring/internal/service"
	"avito-test2024-spring/pkg/auth"
	"avito-test2024-spring/pkg/cache"
	"avito-test2024-spring/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
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

func (h *Handler) Init(host string, port string) *gin.Engine {
	router := gin.New()
	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/metrics"),
		gin.Recovery(),
	)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	metrics.Init()
	router.GET("/metrics", metrics.PrometheusHandler())

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", host, port)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := httpv1.NewHandler(h.bannersService, h.tagsService, h.featuresService, h.usersService, h.logger, h.tokenManager, h.cache)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
