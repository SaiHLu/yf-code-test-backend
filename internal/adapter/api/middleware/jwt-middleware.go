package middleware

import (
	"codetest/internal/adapter/api/presenter"
	"codetest/internal/adapter/api/util"
	portservice "codetest/internal/port/service"

	"github.com/gin-gonic/gin"
)

func AccessTokenMiddleware(jwtService portservice.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := util.GetJwtTokenFromHeader(ctx)
		if err != nil {
			ctx.JSON(401, presenter.JsonResponseWithoutPagination{
				Success: false,
				Error:   err.Error(),
			})
			ctx.Abort()
			return
		}

		userId, err := jwtService.ValidateAccessToken(accessToken)
		if err != nil {
			ctx.JSON(401, presenter.JsonResponseWithoutPagination{
				Success: false,
				Error:   "Invalid access token",
			})
			ctx.Abort()
			return
		}

		ctx.Set("userId", userId)
		ctx.Next()
	}
}

func RefreshTokenMiddleware(jwtService portservice.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		refreshToken, err := util.GetJwtTokenFromHeader(ctx)
		if err != nil {
			ctx.JSON(401, presenter.JsonResponseWithoutPagination{
				Success: false,
				Error:   err.Error(),
			})
			ctx.Abort()
			return
		}

		_, err = jwtService.ValidateRefreshToken(refreshToken)
		if err != nil {
			ctx.JSON(401, presenter.JsonResponseWithoutPagination{
				Success: false,
				Error:   "Invalid refresh token",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
