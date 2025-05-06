package cache

import (
"context"
"encoding/json"
"fmt"
"time"

"github.com/go-redis/redis/v8"
"github.com/online-order-system/inventory-service/config"
)

// RedisCache represents a Redis cache client
type RedisCache struct {
client *redis.Client
ttl    time.Duration
}

// NewRedisCache creates a new Redis cache client
func NewRedisCache(cfg *config.Config) (*RedisCache, error) {
client := redis.NewClient(&redis.Options{
Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
Password: cfg.RedisPassword,
DB:       cfg.RedisDB,
})

// Test connection
ctx := context.Background()
_, err := client.Ping(ctx).Result()
if err != nil {
return nil, fmt.Errorf("failed to connect to Redis: %w", err)
}

return &RedisCache{
client: client,
ttl:    time.Duration(cfg.RedisCacheTTL) * time.Second,
}, nil
}

// Close closes the Redis client
func (c *RedisCache) Close() error {
return c.client.Close()
}

// Get retrieves a value from the cache
func (c *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
data, err := c.client.Get(ctx, key).Result()
if err != nil {
return err
}

return json.Unmarshal([]byte(data), value)
}

// Set stores a value in the cache
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
data, err := json.Marshal(value)
if err != nil {
return err
}

return c.client.Set(ctx, key, data, c.ttl).Err()
}

// Delete removes a value from the cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
return c.client.Del(ctx, key).Err()
}

// FlushAll removes all values from the cache
func (c *RedisCache) FlushAll(ctx context.Context) error {
return c.client.FlushAll(ctx).Err()
}
