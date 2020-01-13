package kronk

import (
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/utkonos-dev/kronk/dlm"
	"github.com/utkonos-dev/kronk/scheduler"
)

var (
	ErrJobNotFound = errors.New("job not found")
	ErrExpiredJob  = errors.New("expired job")
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

func (k *Kronk) Start() {
	k.scheduler.Start()
}

func (k *Kronk) AddRegularJob(name, cronTab string, job func()) error {
	jobID, err := k.scheduler.AddFunc(cronTab, k.wrapFunc(name, job))
	if err != nil {
		return err
	}

	k.jobs.Store(name, jobID)

	return err
}

func (k *Kronk) AddOneTimeJob(name string, runAt time.Time, job func()) error {
	now := time.Now()

	if runAt.Before(now) {
		return ErrExpiredJob
	}

	timer := time.NewTimer(runAt.Sub(now))
	go func() {
		<-timer.C
		k.wrapFunc(name, job)()
	}()

	k.jobs.Store(name, timer)

	return nil
}

func (k *Kronk) RemoveJob(name string) error {
	job, ok := k.jobs.Load(name)
	if !ok {
		return ErrJobNotFound
	}

	v := reflect.ValueOf(job)

	// If job is regular, remove it from scheduler.
	// Otherwise, stop timer.
	if v.Kind() == reflect.String {
		err := k.scheduler.Remove(job.(string))
		if err != nil {
			return err
		}

		k.jobs.Delete(job)
	} else {
		job.(*time.Timer).Stop()
	}

	return nil
}

func (k *Kronk) wrapFunc(name string, job func()) func() {
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
