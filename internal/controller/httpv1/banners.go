package httpv1

import "github.com/gin-gonic/gin"

func (h *Handler) initBannersRoutes(api *gin.RouterGroup) {
	banners := api.Group("/") // add auth middleware for admin
	{
		banners.POST("/", h.bannersAdd)
		banners.PATCH("/:id", h.bannersUpdate)
		banners.DELETE("/:id", h.bannersDelete)
		banners.GET("/", h.bannersGetAll)
	}

	userBanner := api.Group("/user_banner") // add auth middleware for user
	{
		userBanner.GET("/", h.getUserBanner)
	}
}

type bannersAddContent struct {
	Title string `json:"title" binding:"required"`
	Text  string `json:"text" binding:"required"`
	URL   string `json:"url" binding:"required"`
}

type bannersAddInput struct {
	Tags     []int             `json:"tags_ids" binding:"required"`
	Feature  int               `json:"feature_id" binding:"required"`
	Content  bannersAddContent `json:"content" binding:"required"`
	IsActive bool              `json:"is_active" binding:"required"`
}

func (h *Handler) bannersAdd(ctx *gin.Context) {
}

type bannersUpdateContent struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
	URL   string `json:"url,omitempty"`
}

type bannersUpdateInput struct {
	Tags     []int                `json:"tags,omitempty"`
	Feature  int                  `json:"feature,omitempty"`
	Content  bannersUpdateContent `json:"content,omitempty"`
	IsActive bool                 `json:"is_active,omitempty"`
}

func (h *Handler) bannersUpdate(ctx *gin.Context) {

}

func (h *Handler) bannersDelete(ctx *gin.Context) {}

func (h *Handler) bannersGetAll(ctx *gin.Context) {}

func (h *Handler) getUserBanner(ctx *gin.Context) {}
