package glock

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

type User struct {
	Id string `json:"id"`
}

func TestSet(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Password: "",
		Addr:     "localhost:6379",
		DB:       1, // use default DB
	})
	q, err := ConnectDistributedQueue(client, "queue", 3*time.Second)
	if err != nil {
		log.Print(err)
	}
	if err := q.Set("test1", &User{"1"}, &User{"2"}, &User{"3"}, &User{"4"}, &User{"1"}); err != nil {
		t.Fail()
	}
	s, err := q.Size("test1")
	if err != nil || s != 5 {
		t.Fail()
	}
	log.Print(s, err)
}

func TestPushAndPop(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Password: "",
		Addr:     "localhost:6379",
		DB:       1, // use default DB
	})
	q, err := ConnectDistributedQueue(client, "queue", 3*time.Second)
	if err != nil {
		log.Print(err)
	}
	if err := q.Set("test1", &User{"1"}, &User{"2"}, &User{"3"}, &User{"4"}, &User{"5"}, &User{"6"}); err != nil {
		log.Print(err)
		t.Fail()
	}

	if err := q.Push("test1", &User{"Tenguyen"}); err != nil {
		log.Print(2222, err)
	}
	out := &User{}
	if err := q.Pop("test1", out); err != nil || out.Id != "Tenguyen" {
		log.Print(out, err)
		t.Fail()
	}
}

func TestUnshiftAndShift(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Password: "",
		Addr:     "localhost:6379",
		DB:       1, // use default DB
	})
	q, err := ConnectDistributedQueue(client, "queue", 3*time.Second)
	if err != nil {
		log.Print(err)
	}
	if err := q.Set("test1", &User{"1"}, &User{"2"}, &User{"3"}, &User{"4"}, &User{"1"}); err != nil {
		log.Print(err)
		t.Fail()
	}

	if err := q.Unshift("test1", &User{"Tenguyen"}); err != nil {
		log.Print(2222, err)
	}
	out := &User{}
	if err := q.Shift("test1", out); err != nil || out.Id != "Tenguyen" {
		log.Print(out, err)
		t.Fail()
	}
}

func TestList(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Password: "",
		Addr:     "localhost:6379",
		DB:       1, // use default DB
	})
	q, err := ConnectDistributedQueue(client, "queue", 3*time.Second)
	if err != nil {
		log.Print(err)
	}
	out := []*User{}
	err = q.List("test1", 0, 3, func(data []string) error {
		for _, item := range data {
			u := &User{}
			if err := json.Unmarshal([]byte(item), u); err != nil {
				log.Print(err)
				continue
			}
			out = append(out, u)
		}
		return nil
	})
	bin, _ := json.Marshal(out)
	log.Print(11111, string(bin), err)
	out = []*User{}
	err = q.List("test1", 4, 7, func(data []string) error {
		for _, item := range data {
			u := &User{}
			if err := json.Unmarshal([]byte(item), u); err != nil {
				log.Print(err)
				continue
			}
			out = append(out, u)
		}
		return nil
	})
	bin2, _ := json.Marshal(out)
	log.Print(2222222, string(bin2), err)
}
