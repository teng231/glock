package glock

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

/*
	Run test: go test -run TestOptimisticRun
	go test -bench BenchmarkOptimistic100t  -benchmem
	go test -bench BenchmarkOptimistic10000t  -benchmem
?*/
func TestOptimisticRun(t *testing.T) {
	ol, err := CreateOptimisticLock("localhost:6379", "", "test_", time.Second)
	if err != nil {
		panic(err)
	}
	if err := ol.Lock("key1"); err != nil {
		log.Print(err)
	}
	if err := ol.Lock("key1"); err != nil {
		log.Print(err)
	}
	if err := ol.Lock("key1"); err != nil {
		log.Print(err)
	}
	if err := ol.Unlock("key1"); err != nil {
		log.Print(err)
	}
	if err := ol.Lock("key1"); err != nil {
		log.Print(err)
	}
}
func BenchmarkOptimistic100t(t *testing.B) {
	ol, err := CreateOptimisticLock("localhost:6379", "", "test_", time.Second)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		if err := ol.Lock(fmt.Sprintf("key_%v", i)); err != nil {
			log.Print(err)
		}
		if err := ol.Unlock(fmt.Sprintf("key_%v", i)); err != nil {
			log.Print(err)
		}
	}
}
func BenchmarkOptimistic10000Miltiples(t *testing.B) {
	ol, err := CreateOptimisticLock("localhost:6379", "", "test_", time.Second)
	if err != nil {
		panic(err)
	}
	c := make(chan string, 1000)
	wg := &sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		go func() {
			for {
				item := <-c
				if err := ol.Lock(item); err != nil {
					log.Print(err)
				}
				if err := ol.Unlock(item); err != nil {
					log.Print(err)
				}
				wg.Done()
			}

		}()
	}
	for i := 0; i < 10000; i++ {
		c <- fmt.Sprintf("key_%v", i)
		wg.Add(1)
	}
	wg.Wait()
}

func BenchmarkOptimistic10000t(t *testing.B) {
	ol, err := CreateOptimisticLock("localhost:6379", "", "test_", time.Second)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10000; i++ {
		if err := ol.Lock(fmt.Sprintf("key_%v", i)); err != nil {
			log.Print(err)
		}
		if err := ol.Unlock(fmt.Sprintf("key_%v", i)); err != nil {
			log.Print(err)
		}
	}
}
