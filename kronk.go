package kronk

import (
	"time"

	"github.com/utkonos-dev/kronk/dlm"
	"github.com/utkonos-dev/kronk/ems"
	"github.com/utkonos-dev/kronk/scheduler"
)

type Kronk struct {
	dlm       dlm.DLM
	ems       ems.MessageSystem
	scheduler scheduler.Scheduler
	logger    Logger

	defaultLockExp time.Duration
}

type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
}

type Config struct {
	DefaultLockExp time.Duration
}

func NewKronk(locker dlm.DLM, ms ems.MessageSystem, scheduler scheduler.Scheduler, logger Logger, cfg Config) *Kronk {
	return &Kronk{
		dlm:            locker,
		ems:            ms,
		scheduler:      scheduler,
		logger:         logger,
		defaultLockExp: cfg.DefaultLockExp,
	}
}

func (k Kronk) Start() {
	k.scheduler.Start()
}

func (k Kronk) AddJob(name, cronTab string, job func()) error {
	_, err := k.scheduler.AddFunc(cronTab, k.wrapFunc(name, job))
	return err
}

func (k Kronk) wrapFunc(name string, job func()) func() {
	return func() {
		ok, err := k.dlm.Lock(name, k.defaultLockExp)
		if err != nil {
			k.logger.Logf("error locking job %s: %s", name, err.Error())
			return
		}

		if !ok {
			return
		}

		go job()

		err = k.dlm.Unlock(name)
		if err != nil {
			k.logger.Logf("error unlocking job %s: %s", name, err.Error())
			return
		}
	}
}
