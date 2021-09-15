package glock

import (
	"log"
	"testing"
	"time"
)

/*
	Run test: go test -run TestDistributeLock
	go test -bench BenchmarkCounterLock100t  -benchmem
	go test -bench BenchmarkCounterLock10000t  -benchmem
?*/
func TestCountDown(t *testing.T) {
	cd, err := CreateCountLock("localhost:6379", "", "test3_", time.Minute)
	if err != nil {
		log.Print(err)
		t.Fail()
	}
	cd.Start("test1", 100, time.Minute)
	cur, _ := cd.DecrBy("test1", 1)
	log.Print(cur)
	if cur != 99 {
		log.Print(cur)
		t.Fail()
	}
}

func TestIncreaseUp(t *testing.T) {
	cd, err := CreateCountLock("localhost:6379", "", "test3_", time.Minute)
	if err != nil {
		log.Print(err)
		t.Fail()
	}
	cd.Start("test3", 100, time.Minute)
	cur, _ := cd.IncrBy("test3", 1)
	log.Print(cur)
	if cur != 101 {
		log.Print(cur)
		t.Fail()
	}
}

func BenchmarkCounterLock100t(t *testing.B) {
	cd, err := CreateCountLock("localhost:6379", "", "test3_", 4*time.Second)
	if err != nil {
		panic(err)
	}
	cd.Start("locktest1", 0, 4*time.Minute)
	for i := 0; i < 100; i++ {
		if _, err := cd.IncrBy("test4", 1); err != nil {
			log.Print(err)
		}
		if _, err := cd.DecrBy("test4", 1); err != nil {
			log.Print(err)
		}
	}
	cur, err := cd.Current("test4")
	if err != nil {
		log.Print(err)
	}
	if cur != 0 {
		t.Fail()
	}
}

func BenchmarkCounterLock10000t(t *testing.B) {
	cd, err := CreateCountLock("localhost:6379", "", "test3_", 4*time.Second)
	if err != nil {
		panic(err)
	}
	cd.Start("locktest1", 0, 4*time.Minute)
	for i := 0; i < 10000; i++ {
		if _, err := cd.IncrBy("test4", 1); err != nil {
			log.Print(err)
		}
		if _, err := cd.DecrBy("test4", 1); err != nil {
			log.Print(err)
		}
	}
	cur, err := cd.Current("test4")
	if err != nil {
		log.Print(err)
	}
	if cur != 0 {
		t.Fail()
	}
}
