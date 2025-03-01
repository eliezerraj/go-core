package redis_cluster

import (
	"testing"
	"crypto/tls"
	"strings"
	"context"
	"strconv"
	redis "github.com/redis/go-redis/v9"
)

func TestRedisCluster_Add(t *testing.T){

	var envCacheCluster	redis.ClusterOptions
	var redisClusterServer RedisClusterServer

	envCacheCluster.Username = ""
	envCacheCluster.Password = ""
	envCacheCluster.Addrs = strings.Split("clustercfg.memdb-arch.vovqz2.memorydb.us-east-2.amazonaws.com:6379", ",") 

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
		t.Errorf("failed to Put : %s", err)
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