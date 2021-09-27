package glock

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

type IDistributedLock interface {
	// Lock is close request in. Lock all request come to the gate
	Lock(key string) (*LockContext, error)
	// Unlock release the gate
	Unlock(lc *LockContext) error
}

type DistributedLock struct {
	sync     *redsync.Redsync
	timelock time.Duration
	prefix   string
}

type LockContext struct {
	mutex *redsync.Mutex
	ctx   context.Context
}

func CreateDistributedLock(addr, pw, prefix string, timelock time.Duration) (*DistributedLock, error) {
	client := redis.NewClient(&redis.Options{
		Password:        pw,
		Addr:            addr,
		MaxRetries:      10,
		MinRetryBackoff: 15 * time.Millisecond,
		MaxRetryBackoff: 1000 * time.Millisecond,
		DialTimeout:     10 * time.Second,
		PoolSize:        1000,
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
	pool := goredis.NewPool(client)

	rs := redsync.New(pool)
	return &DistributedLock{
		sync:     rs,
		timelock: timelock,
		prefix:   prefix,
	}, nil
}

func (d *DistributedLock) Lock(key string) (*LockContext, error) {
	mutex := d.sync.NewMutex(d.prefix + key)
	ctx, cancel := context.WithTimeout(context.Background(), d.timelock)
	go func() {
		<-ctx.Done()
		cancel()
	}()
	if err := mutex.LockContext(ctx); err != nil {
		cancel()
		return nil, err
	}
	return &LockContext{mutex, ctx}, nil
}

func (d *DistributedLock) Unlock(lc *LockContext) error {
	if _, err := lc.mutex.UnlockContext(lc.ctx); err != nil {
		return err
	}
	return nil
}
