package part

import (
	"context"
	"fmt"

	"github.com/anden007/af_dp_clean_core/misc"
	"github.com/anden007/af_dp_clean_core/pkg"

	"github.com/centrifugal/centrifuge-go"
	"github.com/golang-jwt/jwt"
	jsoniter "github.com/json-iterator/go"
	jsontime "github.com/liamylian/jsontime/v2/v2"
	"github.com/spf13/viper"
)

type Channel struct {
	id           string
	appName      string
	sub          *centrifuge.Subscription
	onReceiveMsg func(info *centrifuge.ClientInfo, content *pkg.MsgContent)
}

type MsgCenter struct {
	appName               string
	token_hmac_secret_key string
	client                *centrifuge.Client
	channels              map[string]*Channel
	jsonEncoder           jsoniter.API
	Enable                bool
	logCenter             ILogCenter
}

type IMsgCenter interface {
	OpenChannel(channels []string, onReceiveMsg func(*centrifuge.ClientInfo, *pkg.MsgContent)) (err error)
	CloseChannel(channels []string) (err error)
	SendMsg(channel, msgType, msg string) (err error)
	CreateAppChannelId(channel string) (result string, err error)
	GetConnToken(user string, exp int64) string
	GetSubscriptionToken(channel string, user string, exp int64) (channelId string, token string)
}

func NewMsgCenter(logCenter ILogCenter) IMsgCenter {
	var err error
	instance := &MsgCenter{
		logCenter: logCenter,
	}
	// loadTime := time.Now()
	appName := viper.GetString("app.name")
	ws_endpoint := viper.GetString("message_center.ws_endpoint")
	token_hmac_secret_key := viper.GetString("message_center.token_hmac_secret_key")
	enable := viper.GetBool("message_center.enable")
	if enable {
		instance.Enable = true
		instance.appName = appName
		instance.token_hmac_secret_key = token_hmac_secret_key
		instance.client = centrifuge.NewJsonClient(ws_endpoint, centrifuge.Config{
			Token: instance.GetConnToken(appName, 0),
		})
		err = instance.client.Connect()
		if err != nil {
			panic(fmt.Sprintf("Connect MessageCenter Server has error:%s", err.Error()))
		}
		instance.channels = make(map[string]*Channel)
		instance.jsonEncoder = jsontime.ConfigWithCustomTimeFormat
	}
	// if ENV == ENUM_ENV_DEV {
	// misc.ServiceLoadInfo("MessageCenter", instance.Enable, loadTime)
	// }
	return instance
}

func (m *MsgCenter) checkMe() {
	if !m.Enable {
		panic("MessageCenter is Disabled.\n")
	}
}

func (m *MsgCenter) OpenChannel(channels []string, onReceiveMsg func(*centrifuge.ClientInfo, *pkg.MsgContent)) (err error) {
	m.checkMe()
	if m.client.State() != centrifuge.StateConnected {
		m.logCenter.Log().Info("[MsgCenter] 尚未连接消息中心, 无法订阅频道")
	} else {
		for _, channel := range channels {
			appChannelId, _ := m.CreateAppChannelId(channel)
			m.logCenter.Log().Infof("[MsgCenter] 准备订阅频道: %s(%s)", channel, appChannelId)
			if _, exists := m.channels[appChannelId]; !exists {
				if sub, sErr := m.client.NewSubscription(appChannelId, centrifuge.SubscriptionConfig{
					GetToken: func(e centrifuge.SubscriptionTokenEvent) (string, error) {
						channelId, token := m.GetSubscriptionToken(channel, m.appName, 0)
						m.logCenter.Log().Infof("[MsgCenter] 生成ChannelToken, channelId: %s token:%s", channelId, token)
						return token, nil
					},
				}); sErr == nil {
					sub.OnSubscribed(func(e centrifuge.SubscribedEvent) {
						m.logCenter.Log().Infof("[MsgCenter] 频道%s, 订阅成功", sub.Channel)
					})
					sub.OnError(func(e centrifuge.SubscriptionErrorEvent) {
						m.logCenter.Log().Infof("[MsgCenter] 频道%s, 发生错误:%s", sub.Channel, e.Error)
					})
					sub.OnUnsubscribed(func(e centrifuge.UnsubscribedEvent) {
						m.logCenter.Log().Infof("[MsgCenter] 频道%s, 取消订阅", sub.Channel)
					})
					sub.OnPublication(func(e centrifuge.PublicationEvent) {
						m.logCenter.Log().Infof("[MsgCenter] 收到消息, 内容: %s", string(e.Data))
						var content *pkg.MsgContent
						err := m.jsonEncoder.Unmarshal(e.Data, &content)
						if err == nil && content != nil {
							m.channels[sub.Channel].onReceiveMsg(e.Info, content)
						}
					})
					err = sub.Subscribe()
					if err == nil {
						channel := Channel{
							id:           channel,
							appName:      m.appName,
							sub:          sub,
							onReceiveMsg: onReceiveMsg,
						}
						m.channels[appChannelId] = &channel
						m.logCenter.Log().Infof("[MsgCenter] 订阅频道%s成功", appChannelId)
					} else {
						m.logCenter.Log().Infof("[MsgCenter] 订阅频道%s发生错误: %s", appChannelId, err.Error())
					}
				}
			}
		}
	}
	return
}

func (m *MsgCenter) CloseChannel(channels []string) (err error) {
	m.checkMe()
	for _, channel := range channels {
		appChannelId, _ := m.CreateAppChannelId(channel)
		sub := m.channels[appChannelId].sub
		if sub != nil {
			err = sub.Unsubscribe()
		}
		delete(m.channels, appChannelId)
	}
	return
}

func (m *MsgCenter) SendMsg(channel, msgType, msg string) (err error) {
	m.checkMe()
	if m.client.State() != centrifuge.StateConnected {
		m.logCenter.Log().Info("[MsgCenter] 尚未连接消息中心, 无法推送消息")
	} else {
		appChannelId, _ := m.CreateAppChannelId(channel)
		msgContent := pkg.MsgContent{
			ChannelId: appChannelId,
			MsgType:   msgType,
			Data:      msg,
		}
		data, _ := m.jsonEncoder.Marshal(msgContent)
		_, err = m.client.Publish(context.Background(), appChannelId, data)
	}
	return
}

func (m *MsgCenter) CreateAppChannelId(channel string) (result string, err error) {
	// 不要检查，否则前端会500错误
	// m.checkMe()
	//此方法可结合SHA等加密算法提供更安全的频道Id生成逻辑
	err = nil
	result = misc.SHA1(fmt.Sprintf("%s_%s", m.appName, channel))
	return
}

func (m *MsgCenter) GetConnToken(user string, exp int64) string {
	// NOTE that JWT must be generated on backend side of your application!
	// Here we are generating it on client side only for example simplicity.
	claims := jwt.MapClaims{"sub": user}
	if exp > 0 {
		claims["exp"] = exp
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(m.token_hmac_secret_key))
	if err != nil {
		panic(err)
	}
	return token
}

func (m *MsgCenter) GetSubscriptionToken(channel string, user string, exp int64) (channelId string, token string) {
	// NOTE that JWT must be generated on backend side of your application!
	// Here we are generating it on client side only for example simplicity.
	channelId, _ = m.CreateAppChannelId(channel)
	claims := jwt.MapClaims{"channel": channelId, "sub": user}
	if exp > 0 {
		claims["exp"] = exp
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(m.token_hmac_secret_key))
	if err != nil {
		panic(err)
	}
	return channelId, token
}
