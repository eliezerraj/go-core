package redis_cluster

import (
	"context"
	"github.com/rs/zerolog/log"

	redis "github.com/redis/go-redis/v9"
)

var childLogger = log.With().Str("go-core.cache", "redis_cluster").Logger()

type RedisClusterServer struct {
	cache *redis.ClusterClient
}

func (r *RedisClusterServer) NewClusterCache(options *redis.ClusterOptions) *RedisClusterServer {
	childLogger.Debug().Msg("NewClusterCache")

	redisClient := redis.NewClusterClient(options)
	return &RedisClusterServer{
		cache: redisClient,
	}
}

func (r *RedisClusterServer) Ping(ctx context.Context) (*string, error) {
	childLogger.Debug().Msg("Ping")

	status, err := r.cache.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *RedisClusterServer) Get(ctx context.Context, key string) (interface{}, error) {
	childLogger.Debug().Msg("Get")

	res, err := r.cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *RedisClusterServer) SetCount(ctx context.Context, key string, valueReg string, value interface{}) (error) {
	childLogger.Debug().Msg("Count")

	_, err := r.cache.HIncrByFloat(ctx, key, valueReg, value.(float64)).Result()
	if err != nil {
		return err
	}

	//s.cache.PExpire(ctx, key, time.Minute * 1).Result()
	return nil
}

func (r *RedisClusterServer) GetCount(ctx context.Context, key string, valueReg string) (interface{}, error) {
	childLogger.Debug().Msg("GetCount")

	res, err := r.cache.HGet(ctx, key, valueReg).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *RedisClusterServer) Set(ctx context.Context, key string, valueReg interface{}) (bool, error) {
	childLogger.Debug().Msg("AddKey")

	err := r.cache.Set(ctx, key, valueReg , 0).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}
