package dlm

import "time"

type DLM interface {
	Lock(key string, exp time.Duration) (success bool, err error)
	Unlock(key string) error
}
