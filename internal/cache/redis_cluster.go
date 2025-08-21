/**
 * @Author Awen
 * @Date 2025/04/04
 * @Email wengaolng@gmail.com
 **/

package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient implements the Cache interface for Redis
type RedisClusterClient struct {
	client redis.UniversalClient
	prefix string
	ttl    time.Duration
}

// NewRedisClient ..
func NewRedisClusterClient(addrs, prefix string, ttl time.Duration, username, password string) (*RedisClusterClient, error) {

	addrStrs := strings.Split(addrs, ",")

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrStrs,
		Username: username,
		Password: password,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &RedisClusterClient{client: client, prefix: prefix, ttl: ttl}, nil
}

// GetCache retrieves a value from Redis
func (c *RedisClusterClient) GetCache(ctx context.Context, key string) (string, error) {
	key = c.prefix + key
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// SetCache stores a value in Redis
func (c *RedisClusterClient) SetCache(ctx context.Context, key, value string) error {
	key = c.prefix + key
	return c.client.Set(ctx, key, value, c.ttl).Err()
}

// DeleteCache delete a value in Redis
func (c *RedisClusterClient) DeleteCache(ctx context.Context, key string) error {
	key = c.prefix + key
	err := c.client.Del(ctx, key).Err()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("redis delete error: %v", err)
	}
	return nil
}

// Close ..
func (c *RedisClusterClient) Close() error {
	return c.client.Close()
}
