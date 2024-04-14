package httpv1

import (
	"avito-test2024-spring/internal/models"
	"avito-test2024-spring/internal/service"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	adminsControlUser := api.Group("/users")
	{
		adminsControlUser.POST("/", h.addUser)
		adminsControlUser.GET("/", h.getAllUsers)
		adminsControlUser.GET("/:id", h.getUserById)
		adminsControlUser.DELETE("/:id", h.deleteUser)
		adminsControlUser.PATCH("/:id", h.updateUser)
	}
}

type UserInput struct {
	IsAdmin bool `json:"is_admin" binding:"required"`
	TagId   int  `json:"tag_id,omitempty"`
}

// @Summary Creates a new user
// @Description Создание нового пользователя
// @Tags user
// @ID create-user
// @Accept json
// @Param body body UserInput true "User creation request"
// @Success 201 {string} string "Пользователь успешно создан"
// @Failure 400 {object} errorResponse "Invalid data provided"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /users [post]
func (h *Handler) addUser(ctx *gin.Context) {
	var user UserInput
	if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, errResp := h.usersService.AddUser(ctx, service.UserAddInput{
		TagId:   user.TagId,
		IsAdmin: user.IsAdmin,
	})
	if errResp.Status != 0 {
		h.logger.Error(ctx, errResp.Status, errResp.Error)
		newErrorResponse(ctx, errResp.Status, errResp.Error)
		return
	}

	ctx.JSON(http.StatusCreated, map[string]string{"access_token": accessToken})
}

// @Summary Обновление пользователя
// @Description Этот эндпоинт предназначен для обновления пользователя по его идентификатору.
// @Tags user
// @ID update-user
// @Accept json
// @Param id path integer true "Идентификатор пользователя"
// @Param body body UserInput true "Запрос на обновление пользователя"
// @Success 200 {string} string "OK"
// @Failure 400 {object} errorResponse "Некорректные данные"
// @Failure 404 {string} string "Пользователь не найден"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /users/{id} [patch]
func (h *Handler) updateUser(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var userInput UserInput
	if err := json.NewDecoder(ctx.Request.Body).Decode(&userInput); err != nil {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user := models.User{
		Id:      userId,
		TagId:   userInput.TagId,
		IsAdmin: userInput.IsAdmin,
	}

	errResp := h.usersService.UpdateUser(ctx, user)
	if errResp.Status != 0 {
		h.logger.Error(ctx, errResp.Status, errResp.Error)
		newErrorResponse(ctx, errResp.Status, errResp.Error)
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Получение пользователя по идентификатору
// @Description Этот эндпоинт предназначен для получения пользователя по его идентификатору.
// @Tags user
// @ID get-user
// @Produce json
// @Param id path integer true "Идентификатор пользователя"
// @Success 200 {object} models.User "OK"
// @Failure 400 {object} errorResponse "Некорректные данные"
// @Failure 404 {string} string "Пользователь не найден"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /users/{id} [get]
func (h *Handler) getUserById(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, errResp := h.usersService.GetUserById(ctx, userId)
	if errResp.Status != 0 {
		h.logger.Error(ctx, errResp.Status, errResp.Error)
		newErrorResponse(ctx, errResp.Status, errResp.Error)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary Получение всех пользователей с фильтрацией по тегу
// @Tags user
// @Description Этот эндпоинт предназначен для получения всех пользователей с фильтрацией по тегу
// @ID get-users
// @Produce json
// @Param tag_id query integer false "Идентификатор тега"
// @Param limit query integer false "Лимит"
// @Param offset query integer false "Оффсет"
// @Success 200 {array} models.User "OK"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /users [get]
func (h *Handler) getAllUsers(ctx *gin.Context) {
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("limit") != "" {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("limit") == "" {
		limit = 0
	}

	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("offset") != "" {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("offset") == "" {
		offset = 0
	}

	tagId, err := strconv.Atoi(ctx.Query("tag_id"))
	if err != nil && errors.Is(err, strconv.ErrSyntax) && ctx.Query("tag_id") != "" {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if ctx.Query("tag_id") == "" {
		tagId = 0
	}

	users, errResp := h.usersService.GetAllUsers(ctx, tagId, limit, offset)
	if errResp.Status != 0 {
		h.logger.Error(ctx, errResp.Status, errResp.Error)
		newErrorResponse(ctx, errResp.Status, errResp.Error)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// @Summary Удаление пользователя по идентификатору
// @Description Этот эндпоинт предназначен для удаления пользователя по его идентификатору.
// @Tags user
// @ID delete-user
// @Param id path integer true "Идентификатор пользователя"
// @Success 204 {string} string "Пользователь успешно удален"
// @Failure 404 {string} string "Пользователь не найден"
// @Failure 500 {object} errorResponse "Внутренняя ошибка сервера"
// @Router /users/{id} [delete]
func (h *Handler) deleteUser(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Error(ctx, http.StatusBadRequest, err.Error())
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	errResp := h.usersService.DeleteUser(ctx, userId)
	if errResp.Status != 0 {
		h.logger.Error(ctx, errResp.Status, errResp.Error)
		newErrorResponse(ctx, errResp.Status, errResp.Error)
		return
	}

	ctx.Status(http.StatusNoContent)
}
