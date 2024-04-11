package cache

import (
	"avito-test2024-spring/internal/config"
	"avito-test2024-spring/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisCache struct {
	ConnPool *redis.Pool

	CacheTTL           time.Duration
	RetryInterval      time.Duration
	MaxNumberOfRetries int
}

func NewRedisCache(cfg config.RedisConfig) *RedisCache {
	return &RedisCache{
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

// https://pkg.go.dev/github.com/gomodule/redigo/redis#example-Args

func (c *RedisCache) Set(banner models.Banner, tagId int, featureId int) error {
	key := fmt.Sprintf("tag_id:%v:feature_id:%v", tagId, featureId)

	conn := c.ConnPool.Get()
	defer conn.Close()

	jsonBanner, err := json.Marshal(banner)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, jsonBanner)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, int64(c.CacheTTL.Seconds()))
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) Get(tagId int, featureId int) (models.Banner, error) {
	key := fmt.Sprintf("tag_id:%v:feature_id:%v", tagId, featureId)

	conn := c.ConnPool.Get()
	defer conn.Close()

	val, err := conn.Do("GET", key)
	if err != nil {
		return models.Banner{}, err
	}
	if val == nil {
		return models.Banner{}, errors.New("not found")
	}

	var banner models.Banner

	err = json.Unmarshal(val.([]byte), &banner)
	if err != nil {
		return models.Banner{}, err
	}

	return banner, nil
}
