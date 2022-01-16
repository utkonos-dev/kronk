package cron

import (
	"github.com/robfig/cron/v3"
	"github.com/utkonos-dev/kronk/scheduler"
	"strconv"
)

type Cron struct {
	processor *cron.Cron
}

func NewScheduler() scheduler.Scheduler {
	return Cron{
		processor: cron.New(),
	}
}

func (c Cron) Start() {
	c.processor.Start()
}

func (c Cron) AddFunc(spec string, cmd func()) (string, error) {
	entryID, err := c.processor.AddFunc(spec, cmd)
	if err != nil {
		return "", nil
	}

	return strconv.Itoa(int(entryID)), nil
}

func (c Cron) Remove(jobID string) error {
	id, err := strconv.Atoi(jobID)
	if err != nil {
		return err
	}

	c.processor.Remove(cron.EntryID(id))
	return nil
}
