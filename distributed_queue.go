package glock

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type IDistributedQueue interface {
	// Push to the left
	Unshift(key string, values ...interface{}) error
	// Push to the right
	Push(key string, values ...interface{}) error
	// release 1 item right
	Pop(key string, out interface{}) error
	// release 1 item right
	Shift(key string, out interface{}) error
	// List item in list
	List(key string, start, stop int64, f func([]string) error) error
	// Size get current size of list
	Size(key string) (int64, error)
}

type DistributedQueue struct {
	client   *redis.Client
	prefix   string
	timelock time.Duration
}

func ConnectDistributedQueue(client *redis.Client, prefix string, timelock time.Duration) (*DistributedQueue, error) {
	if client == nil {
		return nil, errors.New("not found client")
	}
	if timelock < 0 {
		return nil, errors.New("timelock is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timelock)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &DistributedQueue{
		client:   client,
		prefix:   prefix,
		timelock: timelock,
	}, nil
}

func (q *DistributedQueue) Set(key string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), q.timelock)
	defer cancel()

	if err := q.client.Del(ctx, q.prefix+key).Err(); err != nil {
		return err
	}
	if len(values) > 0 {
		inputs := []string{}
		for _, val := range values {
			bin, _ := json.Marshal(val)
			inputs = append(inputs, string(bin))
		}
		if err := q.client.LPush(ctx, q.prefix+key, inputs).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (q *DistributedQueue) Push(key string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), q.timelock)
	defer cancel()
	inputs := []string{}
	for _, val := range values {
		bin, _ := json.Marshal(val)
		inputs = append(inputs, string(bin))
	}
	if err := q.client.RPush(ctx, q.prefix+key, inputs).Err(); err != nil {
		return err
	}
	return nil
}
func (q *DistributedQueue) Pop(key string, out interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), q.timelock)
	defer cancel()
	val, err := q.client.RPop(ctx, q.prefix+key).Result()
	if err != nil && err == redis.Nil {
		return errors.New(Empty)
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(val), out); err != nil {
		return errors.New(CantParse)
	}
	return nil
}

func (q *DistributedQueue) Unshift(key string, values ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), q.timelock)
	defer cancel()
	inputs := []string{}
	for _, val := range values {
		bin, _ := json.Marshal(val)
		inputs = append(inputs, string(bin))
	}
	if err := q.client.LPush(ctx, q.prefix+key, inputs).Err(); err != nil {
		return err
	}
	return nil
}
func (q *DistributedQueue) Shift(key string, out interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), q.timelock)
	defer cancel()
	val, err := q.client.LPop(ctx, q.prefix+key).Result()
	if err != nil && err == redis.Nil {
		return errors.New(Empty)
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(val), out); err != nil {
		return errors.New(CantParse)
	}
	return nil
}

// List start with page 0
func (q *DistributedQueue) List(key string, start, stop int64, f func([]string) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), q.timelock)
	defer cancel()
	vals, err := q.client.LRange(ctx, q.prefix+key, start, stop).Result()
	if err != nil && err == redis.Nil {
		return errors.New(Empty)
	}
	if err != nil {
		return err
	}
	return f(vals)
}

func (q *DistributedQueue) Size(key string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), q.timelock)
	defer cancel()
	c, err := q.client.LLen(ctx, q.prefix+key).Result()
	if err != nil {
		return 0, err
	}
	return c, nil
}
