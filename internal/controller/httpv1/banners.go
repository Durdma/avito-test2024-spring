package httpv1

import "github.com/gin-gonic/gin"

func (h *Handler) initBannersRoutes(api *gin.RouterGroup) {
	banners := api.Group("/") // add auth middleware for admin
	{
		banners.POST("/", h.bannersPOST)
		banners.PATCH("/:id")
		banners.DELETE("/:id")
		banners.GET("/")
	}

	userBanner := api.Group("/user_banner") // add auth middleware for user
	{
		userBanner.GET("/")
	}
}

type bannersPOSTContent struct {
	Title string `json:"title" binding:"required"`
	Text  string `json:"text" binding:"required"`
	URL   string `json:"url" binding:"required"`
}

type bannersPOSTInput struct {
	Tags     []int              `json:"tags_ids" binding:"required"`
	Feature  int                `json:"feature_id" binding:"required"`
	Content  bannersPOSTContent `json:"content" binding:"required"`
	IsActive bool               `json:"is_active" binding:"required"`
}

func (h *Handler) bannersPOST(ctx *gin.Context) {
}

type bannersPatchContent struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
	URL   string `json:"url,omitempty"`
}

type bannersPatchInput struct {
	Tags     []int               `json:"tags,omitempty"`
	Feature  int                 `json:"feature,omitempty"`
	Content  bannersPatchContent `json:"content,omitempty"`
	IsActive bool                `json:"is_active,omitempty"`
}

func (h *Handler) bannersPATCH(ctx *gin.Context) {

}

func (h *Handler) bannersDELETE(ctx *gin.Context) {}

func (h *Handler) bannersGetAll(ctx *gin.Context) {}

func (h *Handler) getUserBanner(ctx *gin.Context) {}
