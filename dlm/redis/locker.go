package redis

import (
	"time"

	"github.com/utkonos-dev/kronk/dlm"

	"github.com/go-redis/redis/v7"
)

type redLocker struct {
	cache redis.UniversalClient
}

func NewLocker(client redis.UniversalClient) dlm.DLM {
	return &redLocker{cache: client}
}

func (s *redLocker) Lock(key string, exp time.Duration) (success bool, err error) {
	res, err := s.cache.SetNX(key, time.Now().String(), exp).Result()
	if err != nil {
		return false, err
	}
	return res, nil
}

func (s *redLocker) Unlock(key string) error {
	return s.cache.Del(key).Err()
}
