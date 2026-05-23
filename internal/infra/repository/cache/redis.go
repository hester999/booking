package cache

import (
	"booking/internal/domain"
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
)

type Booking struct {
	cache *redis.Client
}

func NewCacheRepo(redisClient *redis.Client) *Booking {
	return &Booking{cache: redisClient}
}

func (c *Booking) GetBooking(ctx context.Context, id string) (*domain.Booking, error) {
	redisCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	res, err := c.cache.Get(redisCtx, id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var cache *CacheBooking
	err = json.Unmarshal([]byte(res), &cache)
	if err != nil {
		return nil, err
	}

	restoreBooking, err := cache.CacheToDomain()
	if err != nil {
		return nil, err
	}
	return restoreBooking, nil
}

func (c *Booking) SetBooking(ctx context.Context, booking *domain.Booking) error {
	cacheBooking := &CacheBooking{}

	cacheBooking.DomainToDTOCreate(booking)

	payload, err := json.Marshal(cacheBooking)
	if err != nil {
		return err
	}

	err = c.cache.Set(ctx, booking.ID(), string(payload), 0).Err()

	if err != nil {
		return err
	}
	return nil

}

func (c *Booking) DeleteBooking(ctx context.Context, id string) error {
	err := c.cache.Del(ctx, id).Err()
	if err != nil {
		return err
	}
	return nil
}
