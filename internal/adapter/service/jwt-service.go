package service

import (
	"errors"
	"time"

	"codetest/internal/config"
	"codetest/internal/model"
	portservice "codetest/internal/port/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type jwtService struct {
	cfg *config.AppConfig
}

func NewJWTService(cfg *config.AppConfig) portservice.JWTService {
	return &jwtService{cfg: cfg}
}

func (j *jwtService) GenerateAccessToken(user *model.UserModel) (string, error) {
	claims := &JWTClaims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(j.cfg.ACCESS_TOKEN_TTL))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.cfg.ACCESS_TOKEN_KEY))
}

func (j *jwtService) ValidateAccessToken(token string) (string, error) {
	claims := &JWTClaims{}

	validatedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(j.cfg.ACCESS_TOKEN_KEY), nil
	})

	if err != nil {
		return "", err
	}

	if !validatedToken.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}

func (j *jwtService) GenerateRefreshToken(user *model.UserModel) (string, error) {
	claims := &JWTClaims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(j.cfg.REFRESH_TOKEN_TTL))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.cfg.REFRESH_TOKEN_KEY))
}

func (j *jwtService) ValidateRefreshToken(token string) (string, error) {
	claims := &JWTClaims{}

	validatedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(j.cfg.REFRESH_TOKEN_KEY), nil
	})

	if err != nil {
		return "", err
	}

	if !validatedToken.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}

func (j *jwtService) RefreshToken(refreshToken string) (accessTk, refreshTk string, err error) {
	userID, err := j.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	user := &model.UserModel{
		ID: uuid.MustParse(userID),
	}

	aTk, err := j.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	rTk, err := j.GenerateRefreshToken(user)
	errors.Join(err)
	if err != nil {
		return "", "", err
	}

	return aTk, rTk, nil
}
