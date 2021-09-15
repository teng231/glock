package glock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
)

// Copyright (c) 2017 Pavel Pravosud
// https://github.com/rwz/redis-gcra/blob/master/vendor/perform_gcra_ratelimit.lua

type Limiter struct {
	client   *redis.Client
	timelock time.Duration
	limiter  *redis_rate.Limiter
}

const (
	Second     = "second"
	Minute     = "minute"
	Hour       = "hour"
	Day        = "day"
	Restricted = "restricted"
)

func CreateLimiter(addr, pw string, timelock time.Duration) (*Limiter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw, // no password set
		DB:       1,  // use default DB
	})
	limiter := redis_rate.NewLimiter(client)
	return &Limiter{client, timelock, limiter}, nil
}

func (r *Limiter) Allow(key string, per string, count int) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timelock)
	defer cancel()
	switch per {
	case Second:
		res, err := r.limiter.Allow(ctx, key, redis_rate.PerSecond(count))
		if err != nil {
			return err
		}
		// log.Print("allowed:", res.Allowed, " remaining:", res.Remaining)
		if res.Allowed == 0 {
			return errors.New(Restricted)
		}
		break
	case Minute:
		res, err := r.limiter.Allow(ctx, key, redis_rate.PerMinute(count))
		if err != nil {
			return err
		}
		// log.Print("allowed:", res.Allowed, " remaining:", res.Remaining)
		if res.Allowed == 0 {
			return errors.New(Restricted)
		}
		break
	case Hour:
		res, err := r.limiter.Allow(ctx, key, redis_rate.PerHour(count))
		if err != nil {
			return err
		}
		// log.Print("allowed:", res.Allowed, " remaining:", res.Remaining)
		if res.Allowed == 0 {
			return errors.New(Restricted)
		}
		break
	case Day:
		return r.AllowInDay(key, count)
	}
	return nil
}

func (r *Limiter) Reset(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timelock)
	defer cancel()
	return r.client.Del(ctx, key).Err()
}

func (r *Limiter) AllowInDay(key string, count int) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timelock)
	defer cancel()
	t := time.Now()
	key = fmt.Sprintf("%s_%v", key, t.Unix()/(3600*24))
	timeRemaining := 3600*24 - (60*60*t.Hour() + 60*t.Minute() + t.Second())
	if err := r.client.SetNX(ctx, key, count, time.Duration(timeRemaining)*time.Second).Err(); err != nil {
		return err
	}
	remain, err := r.client.Decr(ctx, key).Result()
	if err != nil {
		return err
	}
	if remain < 0 {
		return errors.New(Restricted)
	}
	return nil
}
