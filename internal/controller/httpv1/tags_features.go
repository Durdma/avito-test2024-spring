package httpv1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initTagsFeaturesRoutes(api *gin.RouterGroup) {
	tags := api.Group("/tags")
	{
		tags.POST("/", h.addTag)
		tags.DELETE("/:id", h.deleteTag)
		tags.GET("/", h.getAllTags)
	}

	features := api.Group("/features")
	{
		features.POST("/", h.addFeature)
		features.DELETE("/:id", h.deleteFeature)
		features.GET("/", h.getAllFeatures)
	}
}

func (h *Handler) addTag(ctx *gin.Context) {
	err := h.tagsService.AddTag(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // TODO add custom error
		return
	}

	ctx.AbortWithStatus(http.StatusCreated)
}

func (h *Handler) deleteTag(ctx *gin.Context) {}

func (h *Handler) getAllTags(ctx *gin.Context) {}

func (h *Handler) addFeature(ctx *gin.Context) {}

func (h *Handler) deleteFeature(ctx *gin.Context) {}

func (h *Handler) getAllFeatures(ctx *gin.Context) {}
