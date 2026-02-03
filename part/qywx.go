package part

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/anden007/dp_clean_core/misc"
	"github.com/anden007/dp_clean_core/pkg"
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	jsontime "github.com/liamylian/jsontime/v2/v2"
	"github.com/monaco-io/request"
)

type QYWXUserInfo struct {
	UserId   string `json:"userid"`
	UserName string `json:"name"`
	Alias    string `json:"alias"`
	Status   int    `json:"status"`
}

type GetUserInfoResponse struct {
	Success  bool         `json:"success"`
	Message  string       `json:"message"`
	UserInfo QYWXUserInfo `json:"userInfo"`
}

type QywxCacheKeyType int

const (
	QywxSessionCacheKey QywxCacheKeyType = iota
	QywxExtDataCacheKey
)

type IQYWXService interface {
	GetAuthUrl(authTarget, requestId string) (result string)
	GetCacheKey(requestId string, cacheType QywxCacheKeyType) string
	TryAuth(ctx iris.Context, authTarget, authBackUrl string) (result QYWXUserInfo, err error)
	ProcessAuthInfo(ctx iris.Context) (result QYWXUserInfo, err error)
	GetAuthedUserInfo(ctx iris.Context) (result QYWXUserInfo, err error)
}

type QYWXService struct {
	appId        string
	agentId      string
	cache        ICache
	sessionCache ISessionCache
	json         jsoniter.API
	logCenter    ILogCenter
}

func NewQYWXService(cache ICache, sessionCache ISessionCache, logCenter ILogCenter) IQYWXService {
	return &QYWXService{
		appId:        os.Getenv("qywx_app_id"),
		agentId:      os.Getenv("qywx_agent_id"),
		cache:        cache,
		sessionCache: sessionCache,
		json:         jsontime.ConfigWithCustomTimeFormat,
		logCenter:    logCenter,
	}
}

func (m *QYWXService) GetAuthUrl(authTarget, requestId string) (result string) {
	redirect_uri := url.QueryEscape(fmt.Sprintf("http://qywx-ap.afocus.net/auth_callback?agentId=%s&env=%s", m.agentId, authTarget))
	return fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s&agentid=%s#wechat_redirect", m.appId, redirect_uri, requestId, m.agentId)
}

func (m *QYWXService) GetCacheKey(requestId string, cacheType QywxCacheKeyType) string {
	switch cacheType {
	case QywxSessionCacheKey:
		return misc.SHA1(fmt.Sprintf("QYWX_Session_%s", requestId))
	case QywxExtDataCacheKey:
		return misc.SHA1(fmt.Sprintf("QYWX_ExtData_%s", requestId))
	default:
		return ""
	}
}

// 第一步：尝试登录，如果已经登录，直接返回
func (m *QYWXService) TryAuth(ctx iris.Context, authTarget, authBackUrl string) (result QYWXUserInfo, err error) {
	authRequestId := ctx.URLParamTrim("authRequestId")
	if authRequestId == "" {
		authRequestId = misc.NewGuidString()
	} else {
		ctx.SetCookieKV("AuthRequestId", authRequestId, iris.CookieExpires(time.Hour*24))
	}
	sessionCacheKey := m.GetCacheKey(authRequestId, QywxSessionCacheKey)
	requestExtDataCacheKey := m.GetCacheKey(authRequestId, QywxExtDataCacheKey)
	if jsonStr, sErr := m.cache.Get(sessionCacheKey); jsonStr != "" && sErr == nil {
		// 缓存中存在用户信息(用户已登录)
		err = m.json.UnmarshalFromString(jsonStr, &result)
		return
	}
	if authBackUrl != "" {
		extDataStr, _ := m.json.MarshalToString(map[string]interface{}{
			"auth_back_url": authBackUrl,
		})
		m.cache.SetEx(requestExtDataCacheKey, extDataStr, time.Minute*5)
	}
	// 跳转到授权链接
	authUrl := m.GetAuthUrl(authTarget, authRequestId)
	ctx.Redirect(authUrl)
	return
}

// 第二步：处理授权信息，如果有返回链接，则跳转
func (m *QYWXService) ProcessAuthInfo(ctx iris.Context) (result QYWXUserInfo, err error) {
	var rsp GetUserInfoResponse
	requestId := ctx.URLParamTrim("requestId")
	if requestId != "" {
		sessionCacheKey := m.GetCacheKey(requestId, QywxSessionCacheKey)
		requestExtDataCacheKey := m.GetCacheKey(requestId, QywxExtDataCacheKey)
		req := request.Client{
			URL:    fmt.Sprintf("http://qywx-ap.afocus.net/get_userInfo?requestId=%s", requestId),
			Method: "GET",
		}
		resp := req.Send().Scan(&rsp)
		if resp.OK() {
			jsonStr := ""
			result = rsp.UserInfo
			jsonStr, err = m.json.MarshalToString(result)
			m.cache.SetEx(sessionCacheKey, jsonStr, time.Hour*24)
			if exists, eErr := m.cache.Exists(requestExtDataCacheKey); exists && eErr == nil {
				if jsonStr, mErr := m.cache.Get(requestExtDataCacheKey); mErr == nil && jsonStr != "" {
					var extData map[string]interface{}
					m.json.UnmarshalFromString(jsonStr, &extData)
					if _authBackUrl, ok := extData["auth_back_url"]; ok {
						authBackUrl := _authBackUrl.(string)
						if authBackUrl != "" && strings.Contains(authBackUrl, "?") {
							authBackUrl = fmt.Sprintf("%s&authRequestId=%s", authBackUrl, requestId)
						} else {
							authBackUrl = fmt.Sprintf("%s?authRequestId=%s", authBackUrl, requestId)
						}
						ctx.Redirect(authBackUrl)
						return
					}
				}
			} else if eErr != nil {
				err = eErr
			}
		} else {
			err = resp.Error()
		}
	}

	if err != nil {
		m.logCenter.Log().WithField("function", "GetQYWXUserInfo").Errorf("获取企业微信用户信息发生错误：%s", err.Error())
	}
	return
}

func (m *QYWXService) GetAuthedUserInfo(ctx iris.Context) (result QYWXUserInfo, err error) {
	authRequestId := ctx.GetCookie("AuthRequestId")
	if authRequestId != "" {
		sessionCacheKey := m.GetCacheKey(authRequestId, QywxSessionCacheKey)
		if jsonStr, sErr := m.cache.Get(sessionCacheKey); jsonStr != "" && sErr == nil {
			// 缓存中存在用户信息(用户已登录)
			err = m.json.UnmarshalFromString(jsonStr, &result)
			return
		}
	}
	return QYWXUserInfo{}, errors.New("用户未登录")
}

func (m *QYWXService) SendMessage(touser string, content string) (err error) {
	var result pkg.APIErrorResult
	if touser != "" {
		req := request.Client{
			URL:    "http://qywx-ap.afocus.net/send_message",
			Method: "POST",
			JSON: iris.Map{
				"msgType": "markdown",
				"touser":  touser,
				"content": content,
				"agentId": 1000107,
			},
		}
		resp := req.Send().Scan(&result)
		if resp.OK() {
			if jsonStr, mErr := m.json.MarshalToString(result); mErr == nil {
				m.logCenter.Log().WithField("function", "SendMessage").Debug(jsonStr)
			}
			if !result.Success {
				err = errors.New(result.Message)
			}
		} else {
			err = resp.Error()
		}
	} else {
		err = errors.New("touser is empty")
	}

	if err != nil {
		m.logCenter.Log().WithField("function", "SendMessage").Errorf("发送企业微信信息发生错误：%s", err.Error())
	}
	return
}
