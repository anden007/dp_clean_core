package part

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/anden007/af_dp_clean_core/misc"

	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
)

type ISessionCache interface {
	Get(ctx iris.Context, key string) (result string, err error)
	Set(ctx iris.Context, key string, value string, ttl time.Duration) (result string, err error)
	Del(ctx iris.Context, key string) (result bool, err error)
	Exists(ctx iris.Context, key string) (result bool, err error)
	Incr(ctx iris.Context, key string) (result int64, err error)
	GetCacheKey(ctx iris.Context, key string) string
}

type SessionCache struct {
	sessionId string
	cache     ICache
}

func NewSessionCache(cache ICache) ISessionCache {
	instance := new(SessionCache)
	instance.cache = cache
	instance.sessionId = instance.md5(viper.GetString("app.name"))
	return instance
}

func (m *SessionCache) Get(ctx iris.Context, key string) (result string, err error) {
	_key := m.GetCacheKey(ctx, key)
	result, err = m.cache.Get(_key)
	return
}

func (m *SessionCache) Set(ctx iris.Context, key string, value string, ttl time.Duration) (result string, err error) {
	_key := m.GetCacheKey(ctx, key)
	result, err = m.cache.SetEx(_key, value, ttl)
	return
}

func (m *SessionCache) Del(ctx iris.Context, key string) (result bool, err error) {
	_key := m.GetCacheKey(ctx, key)
	result, err = m.cache.Del(_key)
	return
}

func (m *SessionCache) Exists(ctx iris.Context, key string) (result bool, err error) {
	_key := m.GetCacheKey(ctx, key)
	result, err = m.cache.Exists(_key)
	return
}

func (m *SessionCache) Incr(ctx iris.Context, key string) (result int64, err error) {
	_key := m.GetCacheKey(ctx, key)
	result, err = m.cache.Incr(_key)
	return
}

func (m *SessionCache) GetCacheKey(ctx iris.Context, key string) string {
	uid := ctx.GetCookie(m.sessionId)
	if uid == "" {
		uid = misc.NewHexId()
		ctx.SetCookieKV(m.sessionId, uid, iris.CookieExpires(time.Hour*24))
	}
	return m.md5(fmt.Sprintf("%s%s", uid, key))
}

func (m *SessionCache) md5(str string) string {
	result := ""
	h := md5.New()
	_, err := h.Write([]byte(str))
	if err == nil {
		result = hex.EncodeToString(h.Sum(nil))
	}
	return result
}
