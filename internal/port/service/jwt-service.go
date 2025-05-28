package portservice

import (
	"codetest/internal/model"
)

type JWTService interface {
	GenerateAccessToken(user *model.UserModel) (token string, err error)

	ValidateAccessToken(token string) (userId string, err error)

	GenerateRefreshToken(user *model.UserModel) (token string, err error)

	ValidateRefreshToken(token string) (userId string, err error)

	RefreshToken(refreshToken string) (accessTk, refreshTk string, err error)
}
