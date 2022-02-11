package glock

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type IOptimisticLock interface {
	Lock(key string, otherDuration ...time.Duration) error
	Unlock(key string) error
}

type OptimisticLock struct {
	client   *redis.Client
	prefix   string
	timelock time.Duration
}

func StartOptimisticLock(cf *ConnectConfig) (*OptimisticLock, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            cf.RedisAddr,
		Password:        cf.RedisPw,
		MaxRetries:      10,
		MinRetryBackoff: 15 * time.Millisecond,
		MaxRetryBackoff: 1000 * time.Millisecond,
		DialTimeout:     10 * time.Second,
		DB:              cf.RedisDb, // use default DB
	})
	if cf.Timelock < 0 {
		return nil, errors.New("timelock is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), cf.Timelock)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &OptimisticLock{
		client,
		cf.Prefix,
		cf.Timelock,
	}, nil
}

// CreateOptimisticLock deprecated
func CreateOptimisticLock(addr, pw, prefix string, timelock time.Duration) (*OptimisticLock, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        pw,
		MaxRetries:      10,
		MinRetryBackoff: 15 * time.Millisecond,
		MaxRetryBackoff: 1000 * time.Millisecond,
		DialTimeout:     10 * time.Second,
		DB:              1, // use default DB
	})
	if timelock < 0 {
		return nil, errors.New("timelock is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timelock)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &OptimisticLock{
		client,
		prefix,
		timelock,
	}, nil
}

func (ol *OptimisticLock) Lock(key string, otherDuration ...time.Duration) error {
	duration := ol.timelock
	if len(otherDuration) == 1 {
		duration = otherDuration[0]
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	allowed, err := ol.client.SetNX(ctx, ol.prefix+key, "ok", duration).Result()
	if err != nil {
		return err
	}
	if !allowed {
		return errors.New(StatusLocked)
	}
	return nil
}

func (ol *OptimisticLock) Unlock(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := ol.client.Del(ctx, ol.prefix+key).Err()
	return err
}
