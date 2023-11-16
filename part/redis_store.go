package part

import (
	"time"
)

type RedisStore struct {
	cache ICache
}

func NewRedisStore(cache ICache) *RedisStore {
	return &RedisStore{cache: cache}
}

func (m *RedisStore) Get(id string, clear bool) string {
	result, err := m.cache.Get(id)
	if err == nil {
		if clear {
			_, _ = m.cache.Del(id)
		}
	}
	return result
}

func (m *RedisStore) Set(id string, value string) error {
	_, err := m.cache.SetEx(id, value, time.Minute)
	return err
}

func (m *RedisStore) Verify(id, answer string, clear bool) bool {
	result := false
	value, err := m.cache.Get(id)
	if err == nil {
		result = answer == value
		if clear {
			_, _ = m.cache.Del(id)
		}
	}
	return result
}
