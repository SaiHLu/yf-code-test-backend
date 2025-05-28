package handler

import (
	"codetest/internal/adapter/api/dto"
	"codetest/internal/adapter/api/middleware"
	"codetest/internal/adapter/api/presenter"
	portservice "codetest/internal/port/service"

	"github.com/gin-gonic/gin"
)

type UserLogHandler struct {
	router      *gin.RouterGroup
	userService portservice.UserLogService
	jwtService  portservice.JWTService
}

func NewUserLogHandler(router *gin.RouterGroup, userService portservice.UserLogService, jwtService portservice.JWTService) *UserLogHandler {
	handler := &UserLogHandler{
		router:      router,
		userService: userService,
		jwtService:  jwtService,
	}

	handler.registerRoutes()

	return handler
}

func (h *UserLogHandler) registerRoutes() {
	route := h.router.Group("/user-logs")
	{
		route.GET("", middleware.AccessTokenMiddleware(h.jwtService), middleware.ValidationMiddleware(&dto.QueryUserLogRequest{}, middleware.BindForm), h.Find)
	}
}

// Find godoc
// @Summary Get User Logs
// @Description Get a list of user logs
// @Tags UserLogs
// @Accept json
// @Produce json
// @Param page query dto.QueryUserLogRequest false "Query params"
// @Success 200 {object} presenter.JsonResponse{data=[]model.UserLogModel}
// @Failure 400 {object} presenter.JsonResponseWithoutPagination
// @Failure 500 {object} presenter.JsonResponseWithoutPagination
// @Security ApiKeyAuth
// @Router /user-logs [get]
func (h *UserLogHandler) Find(c *gin.Context) {
	val, ok := c.Get("validatedRequest")
	if !ok {
		c.JSON(400, presenter.JsonResponseWithoutPagination{
			Success: false,
			Data:    nil,
			Error:   "Invalid request data",
		})
		return
	}

	request, ok := val.(*dto.QueryUserLogRequest)
	if !ok {
		c.JSON(400, presenter.JsonResponseWithoutPagination{
			Success: false,
			Data:    nil,
			Error:   "Invalid request type",
		})
		return
	}

	request.SetDefaultPagination()

	logs, total, err := h.userService.Find(c, request)
	if err != nil {
		c.JSON(500, presenter.JsonResponseWithoutPagination{
			Success: false,
			Data:    nil,
			Error:   "Failed to retrieve user logs: " + err.Error(),
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

	c.JSON(200, presenter.JsonResponse{
		Success:    true,
		Data:       logs,
		Pagination: pagination,
	})
}
