package redis_cluster

import (
	"testing"
	"crypto/tls"
	//"strings"
	"context"
	//"strconv"
	"time"
	redis "github.com/redis/go-redis/v9"
)

var username = "user-01"
var password = "MyPassword-user-01"

//var username = "cache-user-01"
//var password = "CacheUser01Password123!"

//var username = "cache-user-02"
//var password = "CacheUser02Password123!"

var addr = "master.arch-valkey-02.vovqz2.use2.cache.amazonaws.com:6379"
var _value = "value-01"
//var _key  = username + ":issuer:" + username +"-foo-" + _value
var _key  = username + ":issuer:" + username + "-foo-" + _value

// go test -v -run "TestRedisClientPing"
func TestRedisClientPing(t *testing.T){
	var redisClientCache 	RedisClient
	var optRedisClient		redis.Options

	optRedisClient.Username = username
	optRedisClient.Password = password
	optRedisClient.Addr 	= addr
	optRedisClient.PoolSize =     10              // Maximum number of connections in the pool
	optRedisClient.MinIdleConns = 5               // Minimum number of idle connections
	optRedisClient.PoolTimeout =  5 * time.Second // Timeout for getting a connection from the pool

	if true {
		optRedisClient.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	t.Logf("optRedisClient.Username: %v ", optRedisClient.Username)
	t.Logf("optRedisClient.Password: %v ", optRedisClient.Password)
	t.Logf("optRedisClient.Addr: %v ", optRedisClient.Addr)

	clientCache := redisClientCache.NewRedisClientCache(&optRedisClient)
	_, err := clientCache.Ping(context.Background())
	if err != nil {
		t.Errorf("FAILED to ping redis : %s", err)
	}
}

// go test -v -run "TestRedisClientSet"
func TestRedisClientSet(t *testing.T){

	var redisClientCache 	RedisClient
	var optRedisClient		redis.Options

	optRedisClient.Username = username
	optRedisClient.Password = password
	optRedisClient.Addr 	= addr
	optRedisClient.PoolSize =     10              // Maximum number of connections in the pool
	optRedisClient.MinIdleConns = 5               // Minimum number of idle connections
	optRedisClient.PoolTimeout =  5 * time.Second // Timeout for getting a connection from the pool

	if true {
		optRedisClient.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	t.Logf("optRedisClient.Username: %v ", optRedisClient.Username)
	t.Logf("optRedisClient.Password: %v ", optRedisClient.Password)
	t.Logf("optRedisClient.Addr: %v ", optRedisClient.Addr)

	clientCache := redisClientCache.NewRedisClientCache(&optRedisClient)
	_, err := clientCache.Ping(context.Background())
	if err != nil {
		t.Errorf("FAILED to ping redis : %s", err)
	}

	ttl := 30 * time.Minute
	key := _key
	value := _value
	
	t.Logf("key= %v value= %v", key, value)

	res_bol, err := clientCache.Set(context.Background(), key, value, ttl)
	if err != nil {
		t.Errorf("FAILED !!! to Set : %s", err)
	}
	if (!res_bol) {
		t.Errorf("FAILED !!! to Set (FALSE) %v ", res_bol)
	}
}

// go test -v -run "TestRedisClientGet"
func TestRedisClientGet(t *testing.T){

	var redisClientCache 	RedisClient
	var optRedisClient		redis.Options

	optRedisClient.Username = username
	optRedisClient.Password = password
	optRedisClient.Addr = addr
	optRedisClient.PoolSize =     10              // Maximum number of connections in the pool
	optRedisClient.MinIdleConns = 5               // Minimum number of idle connections
	optRedisClient.PoolTimeout =  5 * time.Second // Timeout for getting a connection from the pool

	if true {
		optRedisClient.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	t.Logf("optRedisClient.Username: %v ", optRedisClient.Username)
	t.Logf("optRedisClient.Password: %v ", optRedisClient.Password)
	t.Logf("optRedisClient.Addr: %v ", optRedisClient.Addr)

	clientCache := redisClientCache.NewRedisClientCache(&optRedisClient)
	_, err := clientCache.Ping(context.Background())
	if err != nil {
		t.Errorf("FAILED to ping redis : %s", err)
	}

	key := _key

	t.Logf("key = %v", key)

	res, err := clientCache.Get(context.Background(), key)
	if err != nil {
		t.Errorf("FAILED !!! key = %v %s", key, err)
	} else {
		t.Logf("SUCCESS !!!! key = %v value = %v", key,res)
	}
}

// go test -v -run "TestRedisClientSetGet"
func TestRedisClientAll(t *testing.T){

	var redisClientCache 	RedisClient
	var optRedisClient		redis.Options

	optRedisClient.Username = username
	optRedisClient.Password = password
	optRedisClient.Addr 	= addr
	optRedisClient.PoolSize =     10              // Maximum number of connections in the pool
	optRedisClient.MinIdleConns = 5               // Minimum number of idle connections
	optRedisClient.PoolTimeout =  5 * time.Second // Timeout for getting a connection from the pool

	if true {
		optRedisClient.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	t.Logf("optRedisClient.Username: %v ", optRedisClient.Username)
	t.Logf("optRedisClient.Password: %v ", optRedisClient.Password)
	t.Logf("optRedisClient.Addr: %v ", optRedisClient.Addr)

	clientCache := redisClientCache.NewRedisClientCache(&optRedisClient)
	_, err := clientCache.Ping(context.Background())
	if err != nil {
		t.Errorf("FAILED to ping redis : %s", err)
	}

	ttl := 30 * time.Minute
	key := "user-03" + ":issuer:" + "acme"
	value := "105"
	
	t.Logf("key:%v value:%v", key, value)

	res_bol, err := clientCache.Set(context.Background(), key, value, ttl)
	if err != nil {
		t.Errorf("FAILED !!! to Set : %s", err)
	}
	if (!res_bol) {
		t.Errorf("FAILED !!! to Set (FALSE) %v ", res_bol)
	}

	res, err := clientCache.Get(context.Background(), key)
	if err != nil {
		t.Errorf("FAILED !!! key : %s", err)
	} else {
		t.Logf("SUCCESS !!!! key: %v value: %v", key,res)
	}
}

/*
func TestRedisCluster(t *testing.T){

	var envCacheCluster	redis.ClusterOptions
	var redisClusterServer RedisClusterServer

	envCacheCluster.Username = "user-02"
	envCacheCluster.Password = "MyCachePassword123!"
	envCacheCluster.Addrs = strings.Split("arch-valkey-02-001.arch-valkey-02.vovqz2.use2.cache.amazonaws.com:6379", ",") 

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
}*/