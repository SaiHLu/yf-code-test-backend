package redis

import (
	"codetest/internal/config"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	singletonRedisInstance *redis.Client
	mu                     sync.Mutex
)

type RedisConnection struct {
	*redis.Client
}

func NewRedisConnection(cfg *config.AppConfig) (*RedisConnection, error) {
	mu.Lock()
	defer mu.Unlock()

	if singletonRedisInstance != nil {
		return &RedisConnection{
			Client: singletonRedisInstance,
		}, nil
	}

	options := &redis.Options{
		Addr:     cfg.REDIS_ADDRESS,
		Password: cfg.REDIS_PASSWORD,
		DB:       cfg.REDIS_DB,
	}

	client := redis.NewClient(options)

	singletonRedisInstance = client

	return &RedisConnection{
		Client: singletonRedisInstance,
	}, nil
}

func (r *RedisConnection) Close() error {
	mu.Lock()
	defer mu.Unlock()

	if singletonRedisInstance != nil {
		log.Println("Closing Redis connection...")
		if err := singletonRedisInstance.Close(); err != nil {
			return err
		}
		singletonRedisInstance = nil
		log.Println("Redis connection closed successfully")
	}
	return nil
}

func (r *RedisConnection) GetRedisInstance() *redis.Client {
	return singletonRedisInstance
}
