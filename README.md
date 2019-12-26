# kronk

![kronk](https://thumbs.gfycat.com/SpecificEqualCony-size_restricted.gif)

**kronk** is dead simple scheduler for modern scalable systems. 

It allows you to manage background tasks in a distributed architecture. 

Already supports [Redis](https://redis.io) and [NATS](https://nats.io) as distributed lock manager and message system.

#### Roadmap
- ☑️ Job scheduler and locker
- Job managing across all instances

So, now you can add jobs and be sure that they will be completed on only one instance. Distributed job managing is in progress.

## Get started
### 1. Install and import package

`go get -u github.com/utkonos-dev/kronk`


### 2. Create Kronk
```
import (
    "github.com/utkonos-dev/kronk"
    redisAdapter "github.com/utkonos-dev/kronk/dlm/redis"
    natsAdapter "github.com/utkonos-dev/kronk/ems/nats"
    "github.com/utkonos-dev/kronk/scheduler/cron"
)
```

```
k := kronk.New(
    natsAdapter.NewMS(natsConn),
    redisAdapter.NewLocker(redisConn),
    cron.NewScheduler(),
    logger,
    kronk.Config{
        DefaultLockExp:     time.Second,
        SyncMessageChannel: "sync-kronk",
    },
)
```

### 3. Start scheduler and sync channel listener

```
if err := k.Start(); err != nil {
    // ...
}
```

### 4. Add job

AddJob can be safely called on all instances, but the job will be performed only by one. If you want to dynamically send jobs to all instances, use SendJob.

```
job := func() {
    fmt.Println("That'll work")
}

err := k.AddJob("kronksays", "* * * * *", job)
if err != nil {
    // ...
}
```

## PR accepted!