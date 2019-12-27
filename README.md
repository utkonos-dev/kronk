# kronk [![GoDoc](https://godoc.org/github.com/utkonos-dev/kronk?status.svg)](https://godoc.org/github.com/utkonos-dev/kronk)

**kronk** is dead simple scheduler for modern scalable systems. 

It allows you to manage background tasks in a distributed architecture. 

Already supports [Redis](https://redis.io) as distributed lock manager.

So, now you can add jobs and be sure that they will be completed on only one instance.

## Get started
### 1. Install and import package

`go get -u github.com/utkonos-dev/kronk`


### 2. Create Kronk
```go
import (
    "github.com/utkonos-dev/kronk"
    redisAdapter "github.com/utkonos-dev/kronk/dlm/redis"
    "github.com/utkonos-dev/kronk/scheduler/cron"
)
```

```go
k := kronk.New(
    redisAdapter.NewLocker(redisConn),
    cron.NewScheduler(),
    logger,
    kronk.Config{
        DefaultLockExp:     time.Second,
    },
)
```

### 3. Start scheduler

```go
k.Start()
```

### 4. Add job

AddJob can be safely called on all instances, but the job will be performed only by one.

```go
job := func() {
    fmt.Println("That'll work")
}

err := k.AddJob("kronksays", "* * * * *", job)
if err != nil {
    // ...
}
```

## PR accepted!