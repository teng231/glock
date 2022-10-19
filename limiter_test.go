package glock

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/golang-module/carbon/v2"
)

func TestAllowhour(t *testing.T) {
	r, err := CreateLimiter("localhost:6379", "", 2*time.Second)
	if err != nil {
		panic(err)
	}
	if err := r.Allow("key2", Hour, 5); err != nil {
		log.Print(err)
	}

}
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
	r, _ := StartLimiter(&ConnectConfig{
		RedisAddr: "localhost:6379",
		Timelock:  2 * time.Second,
	})
	r.Reset("rate:key8")
	if err := r.Allow("key8", Day, 5); err != nil {
		log.Print(err)
	}
	time.Sleep(2 * time.Second)
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

func TestDiff2Day(t *testing.T) {
	now := time.Now()

	weekDay := carbon.Time2Carbon(now).SetWeekStartsAt(carbon.Monday).EndOfWeek()
	hours := carbon.Time2Carbon(now).DiffAbsInHours(weekDay)
	log.Print("weekDay:", weekDay, " hours: ", hours)
}

func TestDiff2Hour(t *testing.T) {
	now := time.Now()
	endOfday := carbon.Time2Carbon(now).EndOfDay()
	secs := carbon.Time2Carbon(now).DiffAbsInSeconds(endOfday)
	log.Print("end of day:", endOfday, " hours: ", secs)
}

func TestWeek(t *testing.T) {
	endOfWeekday := carbon.Now().SetWeekStartsAt(carbon.Monday).EndOfWeek().EndOfDay()
	hours := carbon.Now().DiffAbsInHours(endOfWeekday)
	log.Print(endOfWeekday.Carbon2Time().Unix(), "  ", hours)
}

func TestAllowInWeek(t *testing.T) {
	r, _ := StartLimiter(&ConnectConfig{
		RedisAddr: "localhost:6379",
		Timelock:  3 * time.Second,
	})
	r.Reset("key9")
	if err := r.Allow("key9", Week, 5); err != nil {
		log.Print(err)
	}
	time.Sleep(2 * time.Second)
	if err := r.Allow("key9", Week, 5); err != nil {
		log.Print(err)
	}
	if err := r.Allow("key9", Week, 5); err != nil {
		log.Print(err)
	}
	if err := r.Allow("key9", Week, 5); err != nil {
		log.Print(err)
	}
	if err := r.Allow("key9", Week, 5); err != nil {
		log.Print(err)
	}
	if err := r.Allow("key9", Week, 5); err != nil {
		log.Print(err)
	}
}
