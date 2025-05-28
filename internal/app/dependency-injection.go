package app

import (
	"codetest/internal/adapter/api/handler"
	"codetest/internal/adapter/repository/gorm"
	"codetest/internal/adapter/repository/mongo"
	"codetest/internal/adapter/service"
)

func (s *ServerApp) dependencyInjections() error {
	apiRoute := s.Router.Group("/api")

	s.jwtService = service.NewJWTService(s.Cfg)

	s.userLogRepository = mongo.NewUserLogRepository(s.MongoDBConn.Client, "test", "user_logs")
	s.userLogService = service.NewUserLogService(s.userLogRepository)

	s.userRepository = gorm.NewUserRepository(s.PostgresDBConn.GetDBInstance())
	s.userService = service.NewUserService(s.userRepository)
	s.userHandler = handler.NewUserHandler(apiRoute, s.userService, s.jwtService, s.RedisConn.GetRedisInstance(), s.Cfg.REDIS_USER_LOG_CHANNEL)

	s.healthCheckHandler = handler.NewHealthCheckHandler(apiRoute)
	s.swaggerHandler = handler.NewSwaggerHandler(apiRoute)

	s.authHandler = handler.NewAuthHandler(apiRoute, s.userService, s.jwtService)
	s.userLogHandler = handler.NewUserLogHandler(apiRoute, s.userLogService, s.jwtService)
	return nil
}
