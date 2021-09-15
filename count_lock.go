package glock

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type CountLock struct {
	client   *redis.Client
	timelock time.Duration
	prefix   string
}

func CreateCountLock(redisAddr, redisPw, prefix string, timelock time.Duration) (*CountLock, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            redisAddr,
		Password:        redisPw,
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
	return &CountLock{
		client:   client,
		prefix:   prefix,
		timelock: timelock,
	}, nil
}

// Start: khởi động 1 tiến trình bộ đếm. bắt đầu bằng startWith và thời hạn hiệu lực của bộ đếm là expired
// Start có thể dùng để reset counter về 1 giá trị nào đó.
func (cl *CountLock) Start(key interface{}, startWith int, expired time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), cl.timelock)
	defer cancel()

	rdKey := fmt.Sprintf("%v_%v", cl.prefix, key)
	if err := cl.client.Set(ctx, rdKey, startWith, expired).Err(); err != nil {
		return err
	}
	if expired == -1 {
		if err := cl.client.Persist(ctx, rdKey).Err(); err != nil {
			log.Print(err)
			return errors.New(StatusNotPersist)
		}
	}
	return nil
}

func (cl *CountLock) DecrBy(key interface{}, count int64) (int, error) {
	if count <= 0 {
		return 0, errors.New(StatusInvalidCounter)
	}
	ctx, cancel := context.WithTimeout(context.Background(), cl.timelock)
	defer cancel()

	rdKey := fmt.Sprintf("%s_%v", cl.prefix, key)
	curVal, err := cl.client.DecrBy(ctx, rdKey, count).Result()
	if err != nil {
		return 0, err
	}
	return int(curVal), nil
}

func (cl *CountLock) Current(key interface{}) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cl.timelock)
	defer cancel()

	rdKey := fmt.Sprintf("%s_%v", cl.prefix, key)

	current, err := cl.client.Get(ctx, rdKey).Int()
	if err != nil {
		return 0, err
	}
	return current, nil
}

func (cl *CountLock) IncrBy(key interface{}, count int64) (int, error) {
	if count <= 0 {
		return 0, errors.New(StatusInvalidCounter)
	}
	ctx, cancel := context.WithTimeout(context.Background(), cl.timelock)
	defer cancel()

	rdKey := fmt.Sprintf("%s_%v", cl.prefix, key)
	curVal, err := cl.client.IncrBy(ctx, rdKey, count).Result()
	if err != nil {
		return 0, err
	}
	return int(curVal), nil
}
