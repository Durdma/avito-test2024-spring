package controller

import (
	"avito-test2024-spring/internal/controller/httpv1"
	"avito-test2024-spring/internal/service"
	"avito-test2024-spring/pkg/auth"
	"avito-test2024-spring/pkg/cache"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"net/http"
)

type Handler struct {
	bannersService  service.Banners
	tagsService     service.Tags
	featuresService service.Features
	usersService    service.Users
	logger          zerolog.Logger
	tokenManager    auth.TokenManager
	cache           cache.Cache
}

func NewHandler(bannersService service.Banners, tagsService service.Tags,
	featuresService service.Features, usersService service.Users, logger zerolog.Logger, tokenManager auth.TokenManager,
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

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (h *Handler) Init(host string, port string) *gin.Engine {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/metrics", prometheusHandler())

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
