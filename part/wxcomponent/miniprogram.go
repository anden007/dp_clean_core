package wxcomponent

import (
	"errors"
	"fmt"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/auth"
	miniprogramConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/encryptor"
	"github.com/silenceper/wechat/v2/util"
	"github.com/spf13/viper"
)

type IMiniProgram interface {
	GetOpenId(jsCode string) (result string, err error)
	GetUserInfo(jsCode, encryptedData, iv string) (result *encryptor.PlainData, err error)

	SendSubscribe(msg *Message) (err error)
}

type MiniProgram struct {
	wechat  *wechat.Wechat
	Program *miniprogram.MiniProgram
}

func NewMiniProgram(context *Component) *MiniProgram {
	instance := new(MiniProgram)
	instance.wechat = context.Wechat

	cfg := &miniprogramConfig.Config{
		AppID:     viper.GetString("component.weapp.app_id"),
		AppSecret: viper.GetString("component.weapp.app_secret"),
		Cache:     context.InnerCache,
	}
	instance.Program = context.Wechat.GetMiniProgram(cfg)
	return instance
}

func (m *MiniProgram) GetOpenId(jsCode string) (result string, err error) {
	if jsCode != "" {
		session, err := m.Program.GetAuth().Code2Session(jsCode)
		if err == nil {
			result = session.OpenID
		}
	} else {
		result = ""
		err = errors.New("参数jsCode不能空")
	}
	return
}

// Send 发送订阅消息
func (m *MiniProgram) SendSubscribe(msg *Message) (err error) {
	var accessToken string
	accessToken, err = m.Program.GetAuth().GetAccessToken()

	// 发送订阅消息
	// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.send.html
	subscribeSendURL := "https://api.weixin.qq.com/cgi-bin/message/subscribe/send"

	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", subscribeSendURL, accessToken)
	response, err := util.PostJSON(uri, msg)
	if err != nil {
		return
	}
	return util.DecodeWithCommonError(response, "Send")
}

func (m *MiniProgram) GetUserInfo(jsCode, encryptedData, iv string) (result *encryptor.PlainData, err error) {
	if jsCode != "" {
		var session auth.ResCode2Session
		session, err = m.Program.GetAuth().Code2Session(jsCode)
		if err == nil && session.SessionKey != "" {
			if result, err = m.Program.GetEncryptor().Decrypt(session.SessionKey, encryptedData, iv); err == nil {
				result.OpenID = session.OpenID
				result.UnionID = session.UnionID
			}
		}
	} else {
		result = nil
		err = errors.New("参数jsCode不能空")
	}
	return
}
