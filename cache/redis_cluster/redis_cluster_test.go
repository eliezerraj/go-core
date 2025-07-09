package redis_cluster

import (
	"testing"
	"crypto/tls"
	"strings"
	"context"
	"strconv"
	redis "github.com/redis/go-redis/v9"
)

func TestRedisClient(t *testing.T){

	var redisClientCache 	RedisClient
	var optRedisClient		redis.Options

	optRedisClient.Username = "user-04"
	optRedisClient.Password = "MyCachePassword123!"
	optRedisClient.Addr = "master.arch-valkey-02.vovqz2.use2.cache.amazonaws.com:6379" 

	if true {
		optRedisClient.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	t.Logf("optRedisClient.Username: %v ", optRedisClient.Username)

	clientCache := redisClientCache.NewRedisClientCache(&optRedisClient)
	_, err := clientCache.Ping(context.Background())
	if err != nil {
		t.Errorf("failed to ping redis : %s", err)
	}

	key := "user-04" + ":credit_card:" + "number"
	value := "000.111.555.555"
	
	res_bol, err := clientCache.Set(context.Background(), key, value)
	if err != nil {
		t.Errorf("failed to Set : %s", err)
	}
	if (!res_bol) {
		t.Errorf("failed to Set (FALSE) %v ", res_bol)
	}

	res, err := clientCache.Get(context.Background(), key)
	if err != nil {
		t.Errorf("failed to Get : %s", err)
	}

	t.Logf("success key %v result: %v :", key,res)
}

func TestRedisCluster_Add(t *testing.T){

	var envCacheCluster	redis.ClusterOptions
	var redisClusterServer RedisClusterServer

	envCacheCluster.Username = ""
	envCacheCluster.Password = ""
	envCacheCluster.Addrs = strings.Split("arch-valkey-01.vovqz2.ng.0001.use2.cache.amazonaws.com:6379", ",") 

	if !strings.Contains(envCacheCluster.Addrs[0], "127.0.0.1") {
		envCacheCluster.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	cacheRedis := redisClusterServer.NewClusterCache(&envCacheCluster)
	_, err := cacheRedis.Ping(context.Background())
	if err != nil {
		t.Errorf("failed to ping redis : %s", err)
	}

	key := "key-01:" + "01"
	value := 12345.0

	res_bol, err := cacheRedis.Set(context.Background(), key, value)
	if err != nil {
		t.Errorf("failed to Set : %s", err)
	}
	if (!res_bol) {
		t.Errorf("failed to Set (FALSE) %v ", res_bol)
	}

	res, err := cacheRedis.Get(context.Background(), key)
	if err != nil {
		t.Errorf("failed to Get : %s", err)
	}

	value_assert, err := strconv.ParseFloat(res.(string), 64)
	if err != nil {
		t.Errorf("Error : %v", err)
	}

	if (value == value_assert) {
		t.Logf("success result : %v :", value_assert)
	} else {
		t.Errorf("Error : %v %v", value_assert, value)
	}



}