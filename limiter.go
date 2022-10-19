package glock

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"github.com/golang-module/carbon/v2"
)

// Copyright (c) 2017 Pavel Pravosud
// https://github.com/rwz/redis-gcra/blob/master/vendor/perform_gcra_ratelimit.lua
type ILimiter interface {
	// Allow is access comming request and increase counter
	Allow(key string, per string, count int) error
	// Immediate reset counter
	Reset(key string) error
	// Allow using with duration = day
	AllowInDay(key string, count int) error
	AllowInWeek(key string, count int) error
}
type Limiter struct {
	client   *redis.Client
	timelock time.Duration
	limiter  *redis_rate.Limiter
	tz       string // timezone
}

const (
	Second     = "second"
	Minute     = "minute"
	Hour       = "hour"
	Day        = "day"
	Week       = "week"
	Restricted = "restricted"
)

func StartLimiter(cf *ConnectConfig) (*Limiter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cf.RedisAddr,
		Password: cf.RedisPw, // no password set
		DB:       cf.RedisDb, // use default DB
	})
	limiter := redis_rate.NewLimiter(client)
	if cf.Timezone == "" {
		cf.Timezone = "Asia/Ho_Chi_Minh"
	}
	return &Limiter{client, cf.Timelock, limiter, cf.Timezone}, nil
}

// CreateLimiter deprecated
func CreateLimiter(addr, pw string, timelock time.Duration) (*Limiter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw, // no password set
		DB:       1,  // use default DB
	})
	limiter := redis_rate.NewLimiter(client)
	return &Limiter{client, timelock, limiter, "Asia/Ho_Chi_Minh"}, nil
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
	case Minute:
		res, err := r.limiter.Allow(ctx, key, redis_rate.PerMinute(count))
		if err != nil {
			return err
		}
		// log.Print("allowed:", res.Allowed, " remaining:", res.Remaining)
		if res.Allowed == 0 {
			return errors.New(Restricted)
		}
	case Hour:
		res, err := r.limiter.Allow(ctx, key, redis_rate.PerHour(count))
		if err != nil {
			return err
		}
		// log.Print("allowed:", res.Allowed, " remaining:", res.Remaining)
		if res.Allowed == 0 {
			return errors.New(Restricted)
		}
	case Day:
		return r.AllowInDay(key, count)
	case Week:
		return r.AllowInWeek(key, count)
	}
	return nil
}

func (r *Limiter) Reset(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timelock)
	defer cancel()
	// return r.client.Del(ctx, key).Err()
	return r.limiter.Reset(ctx, key)
}

func (r *Limiter) AllowInDay(key string, count int) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timelock)
	defer cancel()
	day := carbon.Now(r.tz).Carbon2Time().Unix() / 86400
	key = fmt.Sprintf("%s_%d", key, day)
	log.Print(key)
	currentValue, err := r.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		// set time expire 1 day for this key
		if err := r.client.SetNX(ctx, key, count, 86400*time.Second).Err(); err != nil {
			return err
		}
		currentValue, _ = r.client.Get(ctx, key).Int64()
	}
	if currentValue <= 0 {
		return errors.New(Restricted)
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

func (r *Limiter) AllowInWeek(key string, count int) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timelock)
	defer cancel()
	now := time.Now()
	weekDay := carbon.Time2Carbon(now).SetWeekStartsAt(carbon.Monday).EndOfWeek()
	hours := carbon.Time2Carbon(now).DiffAbsInHours(weekDay)
	currentValue, err := r.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		if err := r.client.SetNX(ctx, key, count, time.Duration(hours)*time.Hour).Err(); err != nil {
			return err
		}
		currentValue, _ = r.client.Get(ctx, key).Int64()
	}
	if currentValue <= 0 {
		return errors.New(Restricted)
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
