package kronk

import (
	"errors"
	"sync"
	"time"

	"github.com/utkonos-dev/kronk/dlm"
	"github.com/utkonos-dev/kronk/scheduler"
)

var (
	ErrJobNotFound = errors.New("job not found")
)

type Kronk struct {
	jobs sync.Map

	dlm       dlm.DLM
	scheduler scheduler.Scheduler
	logger    Logger

	defaultLockExp time.Duration
}

type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

type Config struct {
	DefaultLockExp time.Duration
}

func New(locker dlm.DLM, scheduler scheduler.Scheduler, logger Logger, cfg Config) *Kronk {
	return &Kronk{
		dlm:            locker,
		scheduler:      scheduler,
		logger:         logger,
		defaultLockExp: cfg.DefaultLockExp,
	}
}

func (k Kronk) Start() {
	k.scheduler.Start()
}

func (k Kronk) AddJob(name, cronTab string, job func()) error {
	jobID, err := k.scheduler.AddFunc(cronTab, k.wrapFunc(name, job))
	if err != nil {
		return err
	}

	k.jobs.Store(name, jobID)

	return err
}

func (k Kronk) RemoveJob(name string) error {
	jobID, ok := k.jobs.Load(name)
	if !ok {
		return ErrJobNotFound
	}

	err := k.scheduler.Remove(jobID.(string))
	if err != nil {
		return err
	}

	k.jobs.Delete(jobID)

	return nil
}

func (k Kronk) wrapFunc(name string, job func()) func() {
	return func() {
		ok, err := k.dlm.Lock(name, k.defaultLockExp)
		if err != nil {
			k.logger.Printf("error locking job %s: %s", name, err.Error())
			return
		}

		if !ok {
			return
		}

		go job()

		err = k.dlm.Unlock(name)
		if err != nil {
			k.logger.Printf("error unlocking job %s: %s", name, err.Error())
			return
		}
	}
}
