package part

import (
	"net/http"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/openplatform"
	openplatformConfig "github.com/silenceper/wechat/v2/openplatform/config"
	"github.com/spf13/viper"

	"github.com/silenceper/wechat/v2/officialaccount/server"
)

type IOpenPlatform interface {
	GetServer(req *http.Request, writer http.ResponseWriter) (result *server.Server)
}

type OpenPlatform struct {
	wechat   *wechat.Wechat
	Platform *openplatform.OpenPlatform
}

func NewOpenPlatform(context *Component) *OpenPlatform {
	instance := new(OpenPlatform)
	instance.wechat = context.Wechat
	cfg := &openplatformConfig.Config{
		AppID:          viper.GetString("component.mp_app_id"),
		AppSecret:      viper.GetString("component.mp_app_secret"),
		Token:          viper.GetString("component.mp_token"),
		EncodingAESKey: viper.GetString("component.mp_aes_key"),
		Cache:          context.InnerCache,
	}
	instance.Platform = context.Wechat.GetOpenPlatform(cfg)
	return instance
}

func (m *OpenPlatform) GetServer(req *http.Request, writer http.ResponseWriter) (result *server.Server) {
	result = m.Platform.GetServer(req, writer)
	return
}
