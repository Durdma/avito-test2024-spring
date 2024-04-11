package cache

import (
	"avito-test2024-spring/internal/config"
	"github.com/gomodule/redigo/redis"
	"time"
)

type Cache struct {
	ConnPool *redis.Pool

	CacheTTL           time.Duration
	RetryInterval      time.Duration
	MaxNumberOfRetries int
}

func NewRedisCache(cfg config.RedisConfig) *Cache {
	return &Cache{
		ConnPool: &redis.Pool{
			IdleTimeout: 5 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", cfg.Host+":"+cfg.Port, redis.DialDatabase(cfg.DB))
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
		CacheTTL:           cfg.CacheTTL,
		RetryInterval:      cfg.RetryInterval,
		MaxNumberOfRetries: cfg.MaxNumberOfRetries,
	}
}
