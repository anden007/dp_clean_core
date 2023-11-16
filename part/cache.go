package part

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/spf13/viper"
)

type ICache interface {
	GetBackend() (result *redis.Client, err error)
	Get(key string) (result string, err error)
	Set(key string, value string) (result string, err error)
	SetEx(key string, value string, ttl time.Duration) (result string, err error)
	Del(key string) (result bool, err error)
	Exists(key string) (result bool, err error)
	Incr(key string) (result int64, err error)
	Decr(key string) (result int64, err error)
	CreateKey(key string) string
	GetMutex(key string, options ...redsync.Option) (result *redsync.Mutex)
	GetRankList(key string, start, stop int64, desc bool) (data []redis.Z, allCount int64, err error)
	UpdateRankScore(key string, member string, score float64) (err error)
	IncrRankScore(key string, member string, score float64) (err error)
	GetMyRank(key string, member string, desc bool) (rank int64, score float64, err error)
}

type CacheConfig struct {
	AppName  string
	Server   string
	Password string
	DataBase int
	PoolSize int
}

type Cache struct {
	appName string
	ctx     context.Context
	redis   *redis.Client
}

func NewCache() ICache {
	// loadTime := time.Now()
	db := viper.GetInt("redis_cache.db")
	poolSize := viper.GetInt("redis_cache.pool_size")
	instance := &Cache{
		ctx:     context.Background(),
		appName: viper.GetString("app.name"),
		redis: redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis_cache.server"),
			Password: viper.GetString("redis_cache.password"),
			DB:       db,
			PoolSize: poolSize,
		}),
	}
	// if ENV == ENUM_ENV_DEV {
	// misc.ServiceLoadInfo("Cache", true, loadTime)
	// }
	return instance
}

func NewCacheByConfig(config CacheConfig) *Cache {
	instance := &Cache{
		ctx:     context.Background(),
		appName: config.AppName,
		redis: redis.NewClient(&redis.Options{
			Addr:     config.Server,
			Password: config.Password,
			DB:       config.DataBase,
			PoolSize: config.PoolSize,
		}),
	}
	return instance
}

func (m *Cache) GetBackend() (result *redis.Client, err error) {
	return m.redis, nil
}

func (m *Cache) GetMutex(key string, options ...redsync.Option) (result *redsync.Mutex) {
	_key := m.CreateKey(key)
	pool := goredis.NewPool(m.redis)
	lock := redsync.New(pool)
	return lock.NewMutex(_key, options...)
}

func (m *Cache) Get(key string) (result string, err error) {
	_key := m.CreateKey(key)
	result, err = m.redis.Get(m.ctx, _key).Result()
	return
}

func (m *Cache) Set(key string, value string) (result string, err error) {
	_key := m.CreateKey(key)
	result, err = m.redis.Set(m.ctx, _key, value, 0).Result()
	return
}

func (m *Cache) SetEx(key string, value string, ttl time.Duration) (result string, err error) {
	_key := m.CreateKey(key)
	result, err = m.redis.Set(m.ctx, _key, value, ttl).Result()
	return
}

func (m *Cache) Del(key string) (result bool, err error) {
	_key := m.CreateKey(key)
	cmdResult, err := m.redis.Del(m.ctx, _key).Result()
	return cmdResult > 0, err
}

func (m *Cache) Exists(key string) (result bool, err error) {
	_key := m.CreateKey(key)
	cmdResult, err := m.redis.Exists(m.ctx, _key).Result()
	return cmdResult > 0, err
}

func (m *Cache) Incr(key string) (result int64, err error) {
	_key := m.CreateKey(key)
	cmdResult, err := m.redis.Incr(m.ctx, _key).Result()
	return cmdResult, err
}

func (m *Cache) Decr(key string) (result int64, err error) {
	_key := m.CreateKey(key)
	cmdResult, err := m.redis.Decr(m.ctx, _key).Result()
	return cmdResult, err
}

func (m *Cache) Md5(str string) string {
	h := md5.New()
	_, _ = h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func (m *Cache) CreateKey(key string) string {
	return fmt.Sprintf("%s_%s", m.appName, key)
}

func (m *Cache) GetRankList(key string, start, stop int64, desc bool) (data []redis.Z, allCount int64, err error) {
	_key := m.CreateKey(key)
	if zCardResult := m.redis.ZCard(m.ctx, _key); zCardResult.Err() == nil {
		allCount = zCardResult.Val()
		if desc {
			if zRangeResult := m.redis.ZRevRangeWithScores(m.ctx, _key, start, stop); zRangeResult.Err() == nil {
				data = zRangeResult.Val()
			} else {
				err = zRangeResult.Err()
			}
		} else {
			if zRangeResult := m.redis.ZRangeWithScores(m.ctx, _key, start, stop); zRangeResult.Err() == nil {
				data = zRangeResult.Val()
			} else {
				err = zRangeResult.Err()
			}
		}
	} else {
		err = zCardResult.Err()
	}
	return
}

func (m *Cache) UpdateRankScore(key string, member string, score float64) (err error) {
	_key := m.CreateKey(key)
	if zAddResult := m.redis.ZAdd(m.ctx, _key, &redis.Z{Member: member, Score: score}); zAddResult.Err() != nil {
		err = zAddResult.Err()
	}
	return
}

func (m *Cache) IncrRankScore(key string, member string, score float64) (err error) {
	_key := m.CreateKey(key)
	if zIncResult := m.redis.ZIncrBy(m.ctx, _key, score, member); zIncResult.Err() != nil {
		err = zIncResult.Err()
	}
	return
}

func (m *Cache) GetMyRank(key string, member string, desc bool) (rank int64, score float64, err error) {
	_key := m.CreateKey(key)
	if desc {
		if zRankResult := m.redis.ZRevRank(m.ctx, _key, member); zRankResult.Err() == nil {
			rank = zRankResult.Val()
		} else {
			err = zRankResult.Err()
		}
	} else {
		if zRankResult := m.redis.ZRank(m.ctx, _key, member); zRankResult.Err() == nil {
			rank = zRankResult.Val()
		} else {
			err = zRankResult.Err()
		}
	}
	if zScoreResult := m.redis.ZScore(m.ctx, _key, member); zScoreResult.Err() == nil {
		score = zScoreResult.Val()
	} else {
		err = zScoreResult.Err()
	}
	return
}
