package glock

import (
	"sync"
	"testing"
)

const number = 100

func TestKmutex(t *testing.T) {
	wg := sync.WaitGroup{}

	km := CreateKmutexInstance()

	ids := []int{}

	for i := 0; i < number; i++ {
		ids = append(ids, i)
	}

	ii := 0
	for i := 0; i < number*number; i++ {
		wg.Add(1)
		go func(iii int) {
			km.Lock(ids[iii])
			km.Unlock(ids[iii])
			wg.Done()
		}(ii)
		ii++
		if ii == number {
			ii = 0
		}
	}
	wg.Wait()
}

func TestWithLock(t *testing.T) {
	wg := sync.WaitGroup{}
	l := sync.Mutex{}
	km := WithLock(&l)

	ids := []int{}

	for i := 0; i < number; i++ {
		ids = append(ids, i)
	}

	ii := 0
	for i := 0; i < number*number; i++ {
		wg.Add(1)
		go func(iii int) {
			km.Lock(ids[iii])
			km.Unlock(ids[iii])
			wg.Done()
		}(ii)
		ii++
		if ii == number {
			ii = 0
		}
	}
	wg.Wait()
}

func TestLockerInterface(t *testing.T) {
	km := CreateKmutexInstance()

	locker := km.Locker("TEST")

	cond := sync.NewCond(locker)

	if false {
		cond.Wait()
	}
}

func BenchmarkKmutex100t(t *testing.B) {
	km := CreateKmutexInstance()
	for i := 0; i < 100; i++ {
		km.Lock(i)
		km.Unlock(i)
	}
}

func BenchmarkKmutex10000t(t *testing.B) {
	km := CreateKmutexInstance()
	for i := 0; i < 10000; i++ {
		km.Lock(i)
		km.Unlock(i)
	}
}

func BenchmarkKmutex100000tMultiples(t *testing.B) {
	km := CreateKmutexInstance()
	wg := &sync.WaitGroup{}
	c := make(chan int, 20000)
	for i := 0; i < 0; i++ {
		go func() {
			for {
				j := <-c
				// log.Print(j)
				km.Lock(j)
				km.Unlock(j)
				wg.Done()
			}
		}()
	}

	for i := 0; i < 100000; i++ {
		c <- i
		wg.Add(1)
	}
	wg.Wait()
}
