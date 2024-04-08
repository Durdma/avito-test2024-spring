package controller

import (
	"avito-test2024-spring/internal/controller/httpv1"
	"avito-test2024-spring/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
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
	handlerV1 := httpv1.NewHandler(h.bannersService, h.tagsService)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
