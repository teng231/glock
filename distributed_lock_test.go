package glock

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

/*
	Run test: go test -run TestDistributeLock
	go test -bench BenchmarkDistributedLock100t  -benchmem
	go test -bench BenchmarkDistributedLock10000t  -benchmem
?*/
func TestDistributeLock(t *testing.T) {
	dl, err := CreateDistributedLock("localhost:6379", "", "test2_", 4*time.Second)
	if err != nil {
		panic(err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		time.AfterFunc(500*time.Millisecond, func() {
			lctx, err := dl.Lock("test-redsync")
			if err != nil {
				log.Print(err)
			}
			log.Print("locked2.1")
			if err := dl.Unlock(lctx); err != nil {
				panic(err)
			}
			log.Print("done 2.2")
			wg.Done()
		})

	}()
	go func() {
		lctx, err := dl.Lock("test-redsync")
		if err != nil {
			log.Print(err)
			log.Panic(err)
		}
		log.Print("locked3.1")
		time.AfterFunc(2*time.Second, func() {
			if err := dl.Unlock(lctx); err != nil {
				panic(err)
			}
			log.Print("done 3.2")
			wg.Done()
		})

	}()
	go func() {
		time.AfterFunc(900*time.Millisecond, func() {
			lctx, err := dl.Lock("test-redsync")
			if err != nil {
				log.Print(err)
				log.Panic(err)
			}
			log.Print("locked1.1")
			time.AfterFunc(2*time.Second, func() {
				if err := dl.Unlock(lctx); err != nil {
					panic(err)
				}
				log.Print("done 1.2")
				wg.Done()
			})
		})

	}()
	wg.Wait()
}

func BenchmarkDistributedLock100t(t *testing.B) {
	dl, err := CreateDistributedLock("localhost:6379", "", "test2_", 4*time.Second)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		lctx, err := dl.Lock(fmt.Sprintf("test%v", i))
		if err != nil {
			log.Print(err)
		}
		if err := dl.Unlock(lctx); err != nil {
			log.Print(err)
		}
	}
}

func BenchmarkDistributedLock10000t(t *testing.B) {
	dl, err := CreateDistributedLock("localhost:6379", "", "test2_", 4*time.Second)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10000; i++ {
		lctx, err := dl.Lock(fmt.Sprintf("test%v", i))
		if err != nil {
			log.Print(err)
		}
		if err := dl.Unlock(lctx); err != nil {
			log.Print(err)
		}
	}
}
