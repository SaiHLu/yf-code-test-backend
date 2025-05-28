package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"codetest/internal/adapter/api/dto"
	"codetest/internal/adapter/api/middleware"
	"codetest/internal/adapter/api/presenter"
	"codetest/internal/model"
	portservice "codetest/internal/port/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	router         *gin.RouterGroup
	userService    portservice.UserService
	jwtService     portservice.JWTService
	redisClient    *redis.Client
	userLogChannel string
}

func NewUserHandler(router *gin.RouterGroup, userService portservice.UserService, jwtService portservice.JWTService, redisClient *redis.Client, userLogChachan string) *UserHandler {
	handler := &UserHandler{
		router:         router,
		userService:    userService,
		jwtService:     jwtService,
		redisClient:    redisClient,
		userLogChannel: userLogChachan,
	}

	handler.registerRoutes()

	return handler
}

func (h *UserHandler) registerRoutes() {
	route := h.router.Group("/users", middleware.AccessTokenMiddleware(h.jwtService))
	{
		route.GET("", middleware.ValidationMiddleware(dto.QueryUserRequest{}, middleware.BindQuery), h.Find)
		route.GET("/:id", middleware.ValidationMiddleware(dto.UserIDParam{}, middleware.BindUri), h.GetOneByID)
		route.POST("", middleware.ValidationMiddleware(dto.CreateUserRequest{}, middleware.BindJSON), h.Create)
		route.PUT("/:id", middleware.ValidationMiddleware(dto.UserIDParam{}, middleware.BindUri), middleware.ValidationMiddleware(dto.UpdateUserRequest{}, middleware.BindJSON), h.Update)
		route.DELETE("/:id", middleware.ValidationMiddleware(dto.UserIDParam{}, middleware.BindUri), h.Delete)
	}
}

// Get Users godoc
// @Summary Get Users
// @Description Get a list of users
// @Tags Users
// @Accept json
// @Produce json
// @Param page query dto.QueryUserRequest false "Query params"
// @Success 200 {object} presenter.JsonResponse{data=[]model.UserModel}
// @Security ApiKeyAuth
// @Router /users [get]
func (h *UserHandler) Find(c *gin.Context) {
	val, ok := c.Get("validatedRequest")
	if !ok {
		c.JSON(400, presenter.JsonResponse{
			Success: false,
			Data:    nil,
			Error:   "Invalid request data",
		})
		return
	}

	request := val.(*dto.QueryUserRequest)

	request.SetDefaultPagination()

	users, total, err := h.userService.Find(c, request)
	if err != nil {
		c.JSON(500, presenter.JsonResponse{
			Success: false,
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	totalPages := (total + int64(request.PageSize) - 1) / int64(request.PageSize)
	pagination := presenter.Pagination{
		TotalCount: total,
		Page:       request.Page,
		PageSize:   request.PageSize,
		TotalPages: totalPages,
	}

	authID, _ := c.Get("userId")
	bytes, _ := json.Marshal(model.UserLogModel{
		UserID: authID.(string),
		Event:  model.UserLogEventRead,
		Data: map[string]string{
			"full_url": c.Request.URL.String(),
		},
		CreatedAt: time.Now(),
	})

	if err := h.redisClient.Publish(c, h.userLogChannel, bytes).Err(); err != nil {
		log.Printf("Failed to publish user reading event to Redis: %v", err.Error())
	}

	c.JSON(200, presenter.JsonResponse{
		Success:    true,
		Data:       users,
		Error:      "",
		Pagination: pagination,
	})
}

// Get User godoc
// @Summary Get User
// @Description Get a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} presenter.JsonResponseWithoutPagination{data=model.UserModel}
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetOneByID(c *gin.Context) {
	val, ok := c.Get("validatedRequest")
	if !ok {
		c.JSON(400, presenter.JsonResponseWithoutPagination{
			Success: false,
			Data:    nil,
			Error:   "Invalid request data",
		})
		return
	}

	request := val.(*dto.UserIDParam)

	user, err := h.userService.GetOneByID(c, uuid.MustParse(request.ID))
	if err != nil {
		c.JSON(500, presenter.JsonResponseWithoutPagination{
			Success: false,
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	authID, _ := c.Get("userId")
	bytes, _ := json.Marshal(model.UserLogModel{
		UserID: authID.(string),
		Event:  model.UserLogEventRead,
		Data: map[string]interface{}{
			"email": user.Email,
			"name":  user.Name,
			"id":    user.ID,
		},
		CreatedAt: time.Now(),
	})

	if err := h.redisClient.Publish(c, h.userLogChannel, bytes).Err(); err != nil {
		log.Printf("Failed to publish user reading event to Redis: %v", err.Error())
	}

	c.JSON(200, presenter.JsonResponseWithoutPagination{
		Success: true,
		Data:    user,
		Error:   "",
	})
}

// Create User godoc
// @Summary Create User
// @Description Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User data"
// @Success 200 {object} presenter.JsonResponseWithoutPagination{data=model.UserModel}
// @Failure 400 {object} presenter.JsonResponseWithoutPagination
// @Failure 422 {object} presenter.JsonResponseWithoutPagination
// @Failure 500 {object} presenter.JsonResponseWithoutPagination
// @Security ApiKeyAuth
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	val, ok := c.Get("validatedRequest")
	if !ok {
		c.JSON(400, presenter.JsonResponseWithoutPagination{
			Success: false,
			Data:    nil,
			Error:   "Invalid request data",
		})
		return
	}

	request := val.(*dto.CreateUserRequest)

	err := h.userService.Create(c, request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, presenter.JsonResponseWithoutPagination{
			Success: false,
			Data:    nil,
			Error:   "Email already exists",
		})
		return
	}

	authID, _ := c.Get("userId")
	bytes, _ := json.Marshal(model.UserLogModel{
		UserID: authID.(string),
		Event:  model.UserLogEventCreate,
		Data: map[string]interface{}{
			"email": request.Email,
			"name":  request.Name,
		},
		CreatedAt: time.Now(),
	})

	if err := h.redisClient.Publish(c, h.userLogChannel, bytes).Err(); err != nil {
		log.Printf("Failed to publish user creating event to Redis: %v", err.Error())
	}

	c.JSON(200, presenter.JsonResponseWithoutPagination{
		Success: true,
		Data:    nil,
		Message: "User created successfully",
	})
}

// Update User godoc
// @Summary Update User
// @Description Update user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserRequest true "User data"
// @Success 200 {object} presenter.JsonResponseWithoutPagination
// @Failure 400 {object} presenter.JsonResponseWithoutPagination
// @Failure 422 {object} presenter.JsonResponseWithoutPagination
// @Failure 500 {object} presenter.JsonResponseWithoutPagination
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	val, ok := c.Get("validatedRequest")
	if !ok {
		c.JSON(400, presenter.JsonResponseWithoutPagination{
			Success: false,
			Data:    nil,
			Error:   "Invalid request data",
		})
		return
	}

	request := val.(*dto.UpdateUserRequest)
	userId := uuid.MustParse(c.Param("id"))

	err := h.userService.Update(c, userId, request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, presenter.JsonResponseWithoutPagination{
			Success: false,
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	authID, _ := c.Get("userId")
	bytes, _ := json.Marshal(model.UserLogModel{
		UserID: authID.(string),
		Event:  model.UserLogEventUpdate,
		Data: map[string]interface{}{
			"email": request.Email,
			"name":  request.Name,
		},
		CreatedAt: time.Now(),
	})

	if err := h.redisClient.Publish(c, h.userLogChannel, bytes).Err(); err != nil {
		log.Printf("Failed to publish user updating event to Redis: %v", err.Error())
	}

	c.JSON(200, presenter.JsonResponseWithoutPagination{
		Success: true,
		Data:    nil,
		Message: "User updated successfully",
	})
}

// Delete User godoc
// @Summary Delete User
// @Description Delete user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} presenter.JsonResponseWithoutPagination
// @Failure 400 {object} presenter.JsonResponseWithoutPagination
// @Failure 422 {object} presenter.JsonResponseWithoutPagination
// @Failure 500 {object} presenter.JsonResponseWithoutPagination
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	userId := uuid.MustParse(c.Param("id"))

	err := h.userService.DeleteOneByID(c, userId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, presenter.JsonResponseWithoutPagination{
			Success: false,
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	authID, _ := c.Get("userId")
	bytes, _ := json.Marshal(model.UserLogModel{
		UserID:    authID.(string),
		Event:     model.UserLogEventDelete,
		Data:      map[string]interface{}{"id": userId},
		CreatedAt: time.Now(),
	})

	if err := h.redisClient.Publish(c, h.userLogChannel, bytes).Err(); err != nil {
		log.Printf("Failed to publish user deleting event to Redis: %v", err.Error())
	}

	c.JSON(200, presenter.JsonResponseWithoutPagination{
		Success: true,
		Data:    nil,
		Message: "User deleted successfully",
	})
}
