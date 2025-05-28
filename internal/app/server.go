package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"codetest/internal/adapter/api/handler"
	"codetest/internal/config"
	"codetest/internal/model"
	"codetest/internal/persistent/mongo"
	"codetest/internal/persistent/postgres"
	"codetest/internal/persistent/redis"
	portrepository "codetest/internal/port/repository"
	portservice "codetest/internal/port/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type ServerApp struct {
	Router         *gin.Engine
	Cfg            *config.AppConfig
	PostgresDBConn *postgres.DBConnection
	MongoDBConn    *mongo.DBConnection
	RedisConn      *redis.RedisConnection

	server *http.Server
	quit   chan os.Signal

	// Dependency injection
	userService       portservice.UserService
	userRepository    portrepository.UserRepository
	jwtService        portservice.JWTService
	userLogService    portservice.UserLogService
	userLogRepository portrepository.UserLogRepository

	swaggerHandler     *handler.SwaggerHandler
	healthCheckHandler *handler.HealthCheckHandler
	userHandler        *handler.UserHandler
	authHandler        *handler.AuthHandler
	userLogHandler     *handler.UserLogHandler
}

func NewServerApp(cfg *config.AppConfig) (*ServerApp, error) {
	gin.SetMode(cfg.MODE)

	engine := gin.New()

	server := &http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: engine.Handler(),
	}

	postgresDBInstance, err := postgres.NewPostgresDBConnection(cfg)
	if err != nil {
		return nil, err
	}

	mongoDBInstance, err := mongo.NewMongoDBConnection(cfg.MONGODB_URI)
	if err != nil {
		return nil, err
	}

	redisInstance, err := redis.NewRedisConnection(cfg)
	if err != nil {
		return nil, err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	return &ServerApp{
		server:         server,
		quit:           quit,
		PostgresDBConn: postgresDBInstance,
		MongoDBConn:    mongoDBInstance,
		RedisConn:      redisInstance,

		Router: engine,
		Cfg:    cfg,
	}, nil
}

func (s *ServerApp) Run() {
	s.middlewares()

	if err := s.dependencyInjections(); err != nil {
		log.Fatalf("Failed to inject dependencies: %v", err)
	}

	if err := s.init(); err != nil {
		panic(err)
	}
}

func (s *ServerApp) init() error {
	errg, errgCtx := errgroup.WithContext(context.Background())

	errg.Go(func() error {
		log.Printf("HTTP server listening on %s\n", s.server.Addr)

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	errg.Go(func() error {
		s.subscribeToUserLogChannel(errgCtx)
		return nil
	})

	errg.Go(func() error {
		<-s.quit
		log.Println("Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(errgCtx, 30*time.Second)
		defer shutdownCancel()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			return err
		}

		if err := s.PostgresDBConn.Close(); err != nil {
			return err
		}

		if err := s.MongoDBConn.Close(shutdownCtx); err != nil {
			return err
		}

		if err := s.RedisConn.Close(); err != nil {
			return err
		}

		log.Println("Server gracefully stopped")
		return nil
	})

	return errg.Wait()
}

func (s *ServerApp) middlewares() {
	s.Router.Use(gin.Recovery())
	s.Router.Use(gin.Logger())

	corsConfig := cors.Config{
		AllowOrigins:     strings.Split(s.Cfg.CORS_ALLOWED_ORIGINS, ","),
		AllowMethods:     strings.Split(s.Cfg.CORS_ALLOWED_METHODS, ","),
		AllowHeaders:     strings.Split(s.Cfg.CORS_ALLOWED_HEADERS, ","),
		ExposeHeaders:    strings.Split(s.Cfg.CORS_EXPOSED_HEADERS, ","),
		AllowCredentials: s.Cfg.CORS_ALLOW_CREDENTIALS,
		MaxAge:           12 * time.Hour,
	}
	s.Router.Use(cors.New(corsConfig))
}

func (s *ServerApp) subscribeToUserLogChannel(ctx context.Context) {
	subscriber := s.RedisConn.Subscribe(ctx, s.Cfg.REDIS_USER_LOG_CHANNEL)
	defer func() {
		_ = subscriber.Close()
		log.Println("Redis subscriber closed")
	}()

	ch := subscriber.Channel()

	for {
		select {
		case <-ctx.Done():
			log.Println("Redis subscription context cancelled")
			return
		case msg, ok := <-ch:
			if !ok {
				log.Println("Redis channel closed")
				return
			}

			var userLogData model.UserLogModel
			if err := json.Unmarshal([]byte(msg.Payload), &userLogData); err != nil {
				log.Printf("Failed to unmarshal user log data: %v", err)
				continue
			}

			if err := s.userLogService.Create(ctx, &userLogData); err != nil {
				log.Printf("Failed to save user log data: %v", err)
				continue
			}
		}
	}
}
