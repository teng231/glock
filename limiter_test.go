package glock

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestAllowSec(t *testing.T) {
	r, err := CreateLimiter("localhost:6379", "", 2*time.Second)
	if err != nil {
		panic(err)
	}
	if err := r.Allow("key1", Second, 5); err != nil {
		log.Print(err)
	}
	if err := r.Allow("key1", Second, 5); err != nil {
		log.Print(err)
	}
	if err := r.Allow("key1", Second, 5); err != nil {
		log.Print(err)
	}
	if err := r.Allow("key1", Second, 5); err != nil {
		log.Print(err)
	}
	if err := r.Allow("key1", Second, 5); err != nil {
		log.Print(err)
	}
	if err := r.Allow("key1", Second, 5); err != nil {
		log.Print(err)
	}
	time.Sleep(1 * time.Second)
	if err := r.Allow("key1", Second, 5); err != nil {
		log.Print(err)
	}
}

func BenchmarkAllow100t(t *testing.B) {
	r, err := CreateLimiter("localhost:6379", "", 2*time.Second)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		if err := r.Allow("key2", Second, 5); err != nil {
			log.Print(err)
		}
	}
}

func BenchmarkAllow10000t(t *testing.B) {
	r, err := CreateLimiter("localhost:6379", "", 2*time.Second)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		if err := r.Allow("key2", Second, 5); err != nil {
			log.Print(err)
		}
	}
}

func BenchmarkAllow10000tMiltiple(t *testing.B) {
	r, err := CreateLimiter("localhost:6379", "", 2*time.Second)
	if err != nil {
		panic(err)
	}
	buf := make(chan int, 20)
	wg := &sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		go func() {
			for {
				<-buf
				if err := r.Allow("key2", Second, 5); err != nil {
					log.Print(err)
				}
				wg.Done()
			}
		}()
	}
	for i := 0; i < 10000; i++ {
		buf <- i
		wg.Add(1)
	}
	wg.Wait()
}

func TestAllowInDay(t *testing.T) {
	r, err := CreateLimiter("localhost:6379", "", 2*time.Second)
	if err != nil {
		panic(err)
	}
	r.Reset("rate:key8")
	if err := r.Allow("key8", Day, 5); err != nil {
		log.Print(err)
	}
	if err := r.Allow("key8", Day, 5); err != nil {
		log.Print(err)
	}
	if err := r.Allow("key8", Day, 5); err != nil {
		log.Print(err)
	}
	if err := r.Allow("key8", Day, 5); err != nil {
		log.Print(err)
	}
	if err := r.Allow("key8", Day, 5); err != nil {
		log.Print(err)
	}
	if err := r.Allow("key8", Day, 5); err != nil {
		log.Print(err)
	}
}
