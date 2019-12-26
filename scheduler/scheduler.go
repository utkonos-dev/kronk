package scheduler

type Scheduler interface {
	Start()
	AddFunc(spec string, cmd func()) (string, error)
	Remove(jobID string) error
}
