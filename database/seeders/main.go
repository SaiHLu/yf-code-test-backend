package main

import (
	"codetest/internal/adapter/repository/gorm"
	"codetest/internal/config"
	"codetest/internal/model"
	"codetest/internal/persistent/postgres"
	"context"
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.NewAppConfig()

	pgConn, err := postgres.NewPostgresDBConnection(cfg)
	if err != nil {
		panic(err)
	}

	userRepo := gorm.NewUserRepository(pgConn.GetDBInstance())

	bytesPass, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		user := &model.UserModel{
			Name:     fmt.Sprintf("User %d", i),
			Email:    fmt.Sprintf("user%d@gmail.com", i),
			Password: string(bytesPass),
		}

		userExists, _ := userRepo.GetOneBy(context.Background(), "email", user.Email)
		if userExists != nil {
			continue
		}

		userRepo.Create(context.Background(), user)
	}
}
