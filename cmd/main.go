package main

import (
	"codetest/internal/app"
	"codetest/internal/config"
	"log"

	_ "github.com/joho/godotenv/autoload"
)

// @title           Yoma Fleet API
// @version         1.0
// @description     This is a sample user management API server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345"
func main() {
	cfg := config.NewAppConfig()

	server, err := app.NewServerApp(cfg)
	if err != nil {
		log.Fatalf("Failed to create server app: %v", err)
	}

	server.Run()
}
