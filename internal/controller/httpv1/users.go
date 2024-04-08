package httpv1

import "github.com/gin-gonic/gin"

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/auth/refresh", h.refreshUserSession)

		adminsControlUser := api.Group("")
		{
			adminsControlUser.POST("/", h.addUser)
			adminsControlUser.GET("/", h.getAllUsers)
			adminsControlUser.DELETE("/:id", h.deleteUser)
		}

	}
}

func (h *Handler) addUser(ctx *gin.Context) {}

func (h *Handler) getAllUsers(ctx *gin.Context) {}

func (h *Handler) deleteUser(ctx *gin.Context) {}

type refreshInput struct {
	Token string `json:"token" binding:"required"`
}

func (h *Handler) refreshUserSession(ctx *gin.Context) {}
