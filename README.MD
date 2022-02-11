# Glock
all lock module and ratelimit using redis for controll.

## install

``` bash
go get github.com/teng231/glock
```

## usage
 many structs and methods, not all but i think it's help you for lock.


### counter lock
Counter lock: It's distributed counter, counter can run up or run down. If run up, you can hanlder check and it's not atomic before all node.
Counter down or coundown it's good.
Helpful for count turn.

``` go
func Run(){
	cd, err := StartCountLock(&ConnectConfig{
		RedisAddr:"localhost:6379",
		RedisPw: "",
		Prefix: "test3_",
		Timelock: time.Minute,
		RedisDb: 1,
	})
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
```

### distributed lock
As You know about this. It's can be lock for all node. Request run one by one, depend on a `key`.
If 2 request going at the same time. Fist request come run and when its done. Second request continue process.

``` go
func Run(){
	dl, err := StartDistributedLock(&ConnectConfig{
		RedisAddr:"localhost:6379",
		RedisPw: "",
		Prefix: "test3_",
		Timelock: time.Minute,
		RedisDb: 1,
	})
	if err != nil {
		panic(err)
	}
    lctx, err := dl.Lock("test-redsync")
    if err := dl.Unlock(lctx); err != nil {
				panic(err)
	}
}
```

If many request coming. Distributed lock it not good because many ping request to check redis. Redis can be down then lock system down.

### kmutex

It's local lock like distributed lock but using mutex and waitgroup. You can using it combine with distributed lock.


``` go
func Run(){
	km := CreateKmutexInstance()
	km.Lock("key")
	defer km.Unlock("key")
}
```

### limiter

handler limit like counter lock. but have duration


``` go
func Run(){
    r, err := StartLimiter(&ConnectConfig{
		RedisAddr:"localhost:6379",
		RedisPw: "",
		Timelock: time.Minute,
		RedisDb: 1,
	})
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
```


### optimistic lock

like distributed lock but request 2 come at the same time. It's will return false.
Redis will lower traffics


``` go
func Run(){
    ol, err := StartOptimisticLock(&ConnectConfig{
		RedisAddr:"localhost:6379",
		RedisPw: "",
		Prefix: "test3_",
		Timelock: time.Minute,
		RedisDb: 1,
	}
	)
	if err != nil {
		panic(err)
	}
	if err := ol.Lock("key1"); err != nil {
		log.Print(err)
	}
}
```
