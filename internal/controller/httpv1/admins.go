package httpv1

import "github.com/gin-gonic/gin"

func (h *Handler) initAdminsRoutes(api *gin.RouterGroup) {
	admins := api.Group("/admins")
	{
		admins.POST("/auth/refresh", h.refreshAdminSession)

		authenticated := api.Group("/")
		{
			authenticated.POST("/", h.addAdmin)
			authenticated.GET("/", h.getAllAdmins)
			authenticated.DELETE("/:id", h.deleteAdmin)
		}
	}
}

func (h *Handler) addAdmin(ctx *gin.Context) {}

func (h *Handler) getAllAdmins(ctx *gin.Context) {}

func (h *Handler) deleteAdmin(ctx *gin.Context) {}

type adminsRefreshInput struct {
	Token string `json:"token" binding:"required"`
}

func (h *Handler) refreshAdminSession(ctx *gin.Context) {}
