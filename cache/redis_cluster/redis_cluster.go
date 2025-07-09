package redis_cluster

import (
	"context"
	"github.com/rs/zerolog/log"

	redis "github.com/redis/go-redis/v9"
)

var childLogger = log.With().Str("component","go-core").Str("package", "cache.redis_cluster").Logger()

type RedisClusterServer struct {
	cache *redis.ClusterClient
}

func (r *RedisClusterServer) NewClusterCache(options *redis.ClusterOptions) *RedisClusterServer {
	childLogger.Debug().Str("func","NewClusterCache").Send()

	redisClient := redis.NewClusterClient(options)
	return &RedisClusterServer{
		cache: redisClient,
	}
}

func (r *RedisClusterServer) Ping(ctx context.Context) (*string, error) {
	childLogger.Debug().Str("func","Ping").Send()

	status, err := r.cache.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *RedisClusterServer) Get(ctx context.Context, key string) (interface{}, error) {
	childLogger.Debug().Str("func","Get").Send()

	res, err := r.cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *RedisClusterServer) SetCount(ctx context.Context, key string, valueReg string, value interface{}) (error) {
	childLogger.Debug().Str("func","SetCount").Send()

	_, err := r.cache.HIncrByFloat(ctx, key, valueReg, value.(float64)).Result()
	if err != nil {
		return err
	}

	//s.cache.PExpire(ctx, key, time.Minute * 1).Result()
	return nil
}

func (r *RedisClusterServer) GetCount(ctx context.Context, key string, valueReg string) (interface{}, error) {
	childLogger.Debug().Str("func","GetCount").Send()

	res, err := r.cache.HGet(ctx, key, valueReg).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *RedisClusterServer) Set(ctx context.Context, key string, valueReg interface{}) (bool, error) {
	childLogger.Debug().Str("func","Set").Send()

	err := r.cache.Set(ctx, key, valueReg , 0).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

type RedisClient struct {
	clientCache *redis.Client
}

func (r *RedisClient) NewRedisClientCache(options *redis.Options) *RedisClient {
	childLogger.Debug().Str("func","NewRedisClientCache").Send()

	redisClient := redis.NewClient(options)

	return &RedisClient{
		clientCache: redisClient,
	}
}

func (r *RedisClient) Ping(ctx context.Context) (*string, error) {
	childLogger.Debug().Str("func","Ping").Send()

	status, err := r.clientCache.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *RedisClient) Set(ctx context.Context, key string, valueReg interface{}) (bool, error) {
	childLogger.Debug().Str("func","Set").Send()

	err := r.clientCache.Set(ctx, key, valueReg , 0).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *RedisClient) Get(ctx context.Context, key string) (interface{}, error) {
	childLogger.Debug().Str("func","Get").Send()

	res, err := r.clientCache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}