package redlock

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
)

type InRedisCache struct {
	pool                     redis.Pool
	maxIdlePoolConnections   uint
	maxActivePoolConnections uint
	maxNumberNamespaces      uint
	databaseNumber           uint
	secondsTTL               uint
}

type PoolOption func(c *InRedisCache)

func WithDatabaseNumber(n uint) PoolOption {
	return func(c *InRedisCache) {
		if n > c.maxNumberNamespaces {
			c.databaseNumber = c.maxNumberNamespaces
			return
		}

		c.databaseNumber = n
	}
}

func WithTTL(seconds uint) PoolOption {
	return func(c *InRedisCache) {
		c.secondsTTL = seconds
	}
}

func NewCache(ctx context.Context, sock string, config ...PoolOption) (*InRedisCache, error) {
	res := InRedisCache{
		maxIdlePoolConnections:   80,
		maxActivePoolConnections: 12000,
		maxNumberNamespaces:      16,
	}

	for _, option := range config {
		option(&res)
	}

	var errConn error

	res.pool = redis.Pool{
		MaxIdle:   res.pool.MaxIdle,
		MaxActive: res.pool.MaxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialContext(ctx, "tcp", sock, redis.DialDatabase(int(res.databaseNumber)))
			if err != nil {
				errConn = err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	if errConn != nil {
		return nil, errConn
	}

	return &res, nil
}
