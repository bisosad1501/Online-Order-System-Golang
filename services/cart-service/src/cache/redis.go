package cache

import (
"context"
"encoding/json"
"fmt"
"time"

"github.com/go-redis/redis/v8"
"github.com/online-order-system/cart-service/config"
)

// RedisCache represents a Redis cache
type RedisCache struct {
client *redis.Client
ttl    time.Duration
}

// NewRedisCache creates a new Redis cache
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
return nil, err
}

return &RedisCache{
client: client,
ttl:    cfg.RedisCacheTTL,
}, nil
}

// Set sets a value in the cache
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
data, err := json.Marshal(value)
if err != nil {
return err
}

return c.client.Set(ctx, key, data, c.ttl).Err()
}

// Get gets a value from the cache
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
data, err := c.client.Get(ctx, key).Bytes()
if err != nil {
return err
}

return json.Unmarshal(data, dest)
}

// Delete deletes a value from the cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
return c.client.Del(ctx, key).Err()
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
return c.client.Close()
}
