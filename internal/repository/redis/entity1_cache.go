package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/LullNil/go-cleanarch/domain"
	"github.com/LullNil/go-cleanarch/domain/entity1"

	goredis "github.com/redis/go-redis/v9"
)

type entity1Cache struct {
	client *goredis.Client
	ttl    time.Duration
}

// NewEntity1Cache creates a Redis cache for entity1.
func NewEntity1Cache(client *goredis.Client, ttl time.Duration) *entity1Cache {
	return &entity1Cache{
		client: client,
		ttl:    ttl,
	}
}

// Get gets an entity1 from cache.
func (c *entity1Cache) Get(ctx context.Context, id int64) (*entity1.Entity1, error) {
	const op = "repository.redis.entity1.Get"

	data, err := c.client.Get(ctx, cacheKey(id)).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var e entity1.Entity1
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &e, nil
}

// Set saves an entity1 to cache.
func (c *entity1Cache) Set(ctx context.Context, e *entity1.Entity1) error {
	const op = "repository.redis.entity1.Set"

	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := c.client.Set(ctx, cacheKey(e.ID), data, c.ttl).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Delete deletes an entity1 from cache.
func (c *entity1Cache) Delete(ctx context.Context, id int64) error {
	const op = "repository.redis.entity1.Delete"

	if err := c.client.Del(ctx, cacheKey(id)).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func cacheKey(id int64) string {
	return fmt.Sprintf("entity1:%d", id)
}
