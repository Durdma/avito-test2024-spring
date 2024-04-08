package controller

import (
	"avito-test2024-spring/internal/controller/httpv1"
	"avito-test2024-spring/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

type Handler struct {
	bannersService service.Banners
	tagsService    service.Tags
	logger         zerolog.Logger
}

func NewHandler(bannersService service.Banners, tagsService service.Tags, logger zerolog.Logger) *Handler {
	return &Handler{
		bannersService: bannersService,
		tagsService:    tagsService,
		logger:         logger,
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

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := httpv1.NewHandler(h.bannersService, h.tagsService, h.logger)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
