package server

import "github.com/go-redis/redis"

var rds *redis.Client

func connectRedis(addr, pass string, db int) error {
	rds = redis.NewClient(&redis.Options{
		Addr: addr,
		Password: pass,
		DB: db,
	})

	err := rds.Ping().Err()
	return err
}