package handler

import (
	"codetest/internal/adapter/api/dto"
	"codetest/internal/adapter/api/middleware"
	"codetest/internal/adapter/api/presenter"
	"codetest/internal/adapter/api/util"
	portservice "codetest/internal/port/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	router      *gin.RouterGroup
	userService portservice.UserService
	jwtService  portservice.JWTService
}

func NewAuthHandler(router *gin.RouterGroup, userService portservice.UserService, jwtService portservice.JWTService) *AuthHandler {
	handler := &AuthHandler{
		router:      router,
		userService: userService,
		jwtService:  jwtService,
	}

	handler.registerRoutes()

	return handler
}

func (h *AuthHandler) registerRoutes() {
	route := h.router.Group("/auth")
	{
		route.POST("/login", middleware.ValidationMiddleware(dto.LoginRequest{}, middleware.BindForm), h.Login)
		route.GET("/me", middleware.AccessTokenMiddleware(h.jwtService), h.Me)
		route.POST("/refresh-token", middleware.RefreshTokenMiddleware(h.jwtService), h.RefreshToken)
	}
}

// Login godoc
// @Summary User Login
// @Description User login to get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} presenter.JsonResponseWithoutPagination{data=dto.LoginResponse}
// @Failure 400 {object} presenter.JsonResponseWithoutPagination
// @Failure 401 {object} presenter.JsonResponseWithoutPagination
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	val, ok := c.Get("validatedRequest")
	if !ok {
		c.JSON(400, presenter.JsonResponseWithoutPagination{
			Success: false,
			Error:   "Invalid request data",
		})
		return
	}

	request, ok := val.(*dto.LoginRequest)
	if !ok {
		c.JSON(400, presenter.JsonResponseWithoutPagination{
			Success: false,
			Error:   "Invalid request type",
		})
		return
	}

	user, err := h.userService.GetOneByEmail(c, request.Email)
	if err != nil {
		c.JSON(400, presenter.JsonResponseWithoutPagination{
			Success: false,
			Error:   "Invalid email or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.JSON(400, presenter.JsonResponseWithoutPagination{
			Success: false,
			Error:   "Invalid email or password",
		})
		return
	}

	accessToken, err := h.jwtService.GenerateAccessToken(user)
	if err != nil {
		c.JSON(400, presenter.JsonResponseWithoutPagination{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(400, presenter.JsonResponseWithoutPagination{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, presenter.JsonResponseWithoutPagination{
		Success: true,
		Data: dto.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// Me godoc
// @Summary Get Current Auth User
// @Description Get current authenticated user information
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} presenter.JsonResponseWithoutPagination{data=model.UserModel}
// @Failure 401 {object} presenter.JsonResponseWithoutPagination
// @Security ApiKeyAuth
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userId, _ := c.Get("userId")

	user, err := h.userService.GetOneByID(c, uuid.MustParse(userId.(string)))
	if err != nil {
		c.JSON(401, presenter.JsonResponseWithoutPagination{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	c.JSON(200, presenter.JsonResponseWithoutPagination{
		Success: true,
		Data:    user,
	})
}

// RefreshToken godoc
// @Summary Refresh JWT Token
// @Description Refresh JWT token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} presenter.JsonResponseWithoutPagination{data=dto.LoginResponse}
// @Failure 401 {object} presenter.JsonResponseWithoutPagination
// @Security ApiKeyAuth
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := util.GetJwtTokenFromHeader(c)
	if err != nil {
		c.JSON(401, presenter.JsonResponseWithoutPagination{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	accessToken, refreshToken, err := h.jwtService.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(401, presenter.JsonResponseWithoutPagination{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(200, presenter.JsonResponseWithoutPagination{
		Success: true,
		Data: dto.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}
