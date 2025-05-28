package util

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetJwtTokenFromHeader(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("invalid Authorization header")
	}

	tokenParts := strings.Split(authHeader, "Bearer ")
	if len(tokenParts) != 2 {
		return "", errors.New("invalid Authorization header")
	}

	return strings.TrimSpace(tokenParts[1]), nil
}
