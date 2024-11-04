package wxcomponent

import (
	"context"
	"time"

	"github.com/anden007/dp_clean_core/misc"
	"github.com/anden007/dp_clean_core/part"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/spf13/viper"
)

type IComponent interface {
	GetMiniProgram() (result IMiniProgram)
	GetOpenPlatform() (result IOpenPlatform)
}

type Component struct {
	Wechat       *wechat.Wechat
	InnerCache   *cache.Redis
	MiniProgram  IMiniProgram
	OpenPlatform IOpenPlatform
}

func NewComponent() *Component {
	instance := new(Component)
	loadTime := time.Now()
	enable := viper.GetBool("component.enable")
	if enable {
		enable = true
		instance.Wechat = wechat.NewWechat()
		weixin_redis_db := viper.GetInt("component.redis_db")
		weixin_redis_pool_size := viper.GetInt("component.redis_pool_size")
		instance.InnerCache = cache.NewRedis(context.Background(), &cache.RedisOpts{
			Host:        viper.GetString("component.redis_server"),
			Password:    viper.GetString("component.redis_password"),
			Database:    weixin_redis_db,
			MaxIdle:     weixin_redis_pool_size / 10,
			MaxActive:   weixin_redis_pool_size,
			IdleTimeout: 1000,
		})
		if instance.MiniProgram == nil {
			instance.MiniProgram = NewMiniProgram(instance)
		}
	}
	if part.ENV == part.ENUM_ENV_DEV {
		misc.ServiceLoadInfo("WxComponent", enable, loadTime)
	}
	return instance
}

func (m *Component) GetMiniProgram() (result IMiniProgram) {
	return m.MiniProgram
}

func (m *Component) GetOpenPlatform() (result IOpenPlatform) {
	return m.OpenPlatform
}
