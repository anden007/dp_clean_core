package part

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/anden007/af_dp_clean_core/misc"
	"github.com/anden007/af_dp_clean_core/pkg"

	"github.com/disintegration/imaging"
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	jsonTime "github.com/liamylian/jsontime/v2/v2"
	"github.com/spf13/viper"
)

type IWxAuth interface {
	GetCurrentUser(ctx iris.Context) (result *pkg.WxAuthInfo)
	CheckAuth(ctx iris.Context, cak string, customReturnURL string) (result bool)
	ClearAuth(ctx iris.Context) (err error)
	DoAuth(ctx iris.Context, afterAuth func(authInfo *pkg.WxAuthInfo, cak string))
	GetAuthUrl(cak string) (result string)
}

type WxAuth struct {
	dev_mock_info        bool
	campaignId           string
	returnUrl            string
	authUrl              string
	authInfoUrl          string
	jwtInstance          IJwtService
	cacheInstance        ICache
	sessionCacheInstance ISessionCache
	json                 jsoniter.API
	logCenter            ILogCenter
	Enable               bool
}

func NewWxAuth(cache ICache, sessionCache ISessionCache, logCenter ILogCenter) IWxAuth {
	instance := new(WxAuth)
	enable := viper.GetString("wxauth.enable")
	if strings.EqualFold("true", enable) {
		instance.Enable = true
		instance.dev_mock_info = viper.GetBool("wxauth.dev_mock_info")
		instance.logCenter = logCenter
		instance.json = jsonTime.ConfigWithCustomTimeFormat
		instance.cacheInstance = cache
		instance.sessionCacheInstance = sessionCache
		instance.campaignId = viper.GetString("wxauth.campaignid")
		instance.returnUrl = viper.GetString("wxauth.return_url")
		instance.authUrl = viper.GetString("wxauth.auth_url")
		instance.authInfoUrl = viper.GetString("wxauth.auth_info_url")
		instance.jwtInstance = NewCustomJWT(fmt.Sprintf("%s_wxtoken", viper.GetString("app.name")))
	}
	return instance
}

func (m *WxAuth) checkMe() {
	if !m.Enable {
		panic("WxAuth is Disabled. if you need enable it,please set 'wxauth_enable = true' in .env file.\n")
	}
}

func (m *WxAuth) GetAuthUrl(cak string) (result string) {
	if strings.Contains(cak, "_") {
		panic("CAK中不能包含'_'，因为它有专门的用途")
	}
	result = fmt.Sprintf("%s%s?cak=%s", m.authUrl, m.campaignId, cak)
	return
}

func (m *WxAuth) getAuthInfoUrl(authId string) (result string) {
	result = fmt.Sprintf("%s%s", m.authInfoUrl, authId)
	return
}

func (m *WxAuth) GetCurrentUser(ctx iris.Context) (result *pkg.WxAuthInfo) {
	m.checkMe()
	if ENV == ENUM_ENV_DEV && m.dev_mock_info {
		result = &pkg.WxAuthInfo{
			ID:         "99999999-9999-9999-9999-999999999999",
			CAK:        "",
			AppID:      "",
			OpenID:     viper.GetString("wxauth_debug_openid"),
			UnionId:    "",
			NickName:   "本地测试用户",
			HeadImgURL: "https://img.afocus.net/public/img/devheadimg.png",
			Sex:        "男",
			Country:    "中国",
			Province:   "重庆",
			City:       "重庆",
			Subscribe:  0,
			CreateTime: time.Now(),
		}
	} else {
		if verifiedToken, err := m.jwtInstance.VerifyToken(ctx); err == nil {
			wxAuthInfoClaims := pkg.WxAuthInfoClaims{}
			if err := verifiedToken.Claims(&wxAuthInfoClaims); err == nil {
				cak := ctx.URLParam("cak")
				result = &wxAuthInfoClaims.WxAuthInfo
				result.CAK = cak
			}
		}
	}
	return
}

func (m *WxAuth) ClearAuth(ctx iris.Context) (err error) {
	m.jwtInstance.RemoveToken(ctx)
	return
}

func (m *WxAuth) CheckAuth(ctx iris.Context, cak string, customReturnURL string) (result bool) {
	m.checkMe()
	if m.GetCurrentUser(ctx) == nil {
		url := m.GetAuthUrl(cak)
		result = false
		if customReturnURL != "" {
			m.sessionCacheInstance.Set(ctx, "auth_return_url", customReturnURL, time.Minute*5)
		}
		ctx.Redirect(url)
	} else {
		result = true
	}
	return
}

func (m *WxAuth) DoAuth(ctx iris.Context, afterAuth func(authInfo *pkg.WxAuthInfo, cak string)) {
	m.checkMe()
	customReturnURL := ""
	if cacheReturnUrl, cacheErr := m.sessionCacheInstance.Get(ctx, "auth_return_url"); cacheErr == nil {
		m.sessionCacheInstance.Del(ctx, "auth_return_url")
		customReturnURL = cacheReturnUrl
	}
	wxAuthInfo := &pkg.WxAuthInfo{}
	authid := ctx.URLParam("authid")
	cak := ctx.URLParam("cak")
	if authid != "" {
		// http模式
		url := m.getAuthInfoUrl(authid)
		if resp, err := http.Get(url); err != nil {
			m.logCenter.Log().Errorf("用户授权失败，请求：%s 错误：%s", url, err.Error())
			_, _ = ctx.Text("用户授权失败，错误：%s", err.Error())
		} else {
			defer resp.Body.Close()
			if body, err := ioutil.ReadAll(resp.Body); err == nil {
				if err := m.json.Unmarshal(body, &wxAuthInfo); err != nil {
					m.logCenter.Log().Errorf("用户信息格式错误：%s", err.Error())
					_, _ = ctx.Text("用户信息格式错误：%s", err.Error())
				}
			}
		}
		if wxAuthInfo != nil {
			// 保存头像到系统（文件中心或本地），避免微信头像过期无法显示
			protocol := viper.GetString("file_center.protocol")
			max_width := 132
			max_height := 0
			image_domain := viper.GetString("file_center.image_domain")
			file_center_url := fmt.Sprintf("%s%s", viper.GetString("file_center.address"), viper.GetString("app.name"))
			if resp, gErr := http.Get(wxAuthInfo.HeadImgURL); gErr == nil {
				defer resp.Body.Close()
				switch protocol {
				case "http", "https":
					if resp, pErr := http.Post(file_center_url, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("max_width=%d&remote_url=%s", 132, wxAuthInfo.HeadImgURL))); pErr == nil {
						defer resp.Body.Close()
						upload_result := pkg.APIResult{}
						if uploadResultContent, rErr := ioutil.ReadAll(resp.Body); rErr == nil {
							if mErr := m.json.Unmarshal(uploadResultContent, &upload_result); mErr == nil {
								wxAuthInfo.HeadImgURL = upload_result.Result.(string)
							}
						}
					}
				case "local":
					fPath := "./uploads/" + time.Now().Format("2006-01")
					fileExt := ".png"
					fName := misc.NewGuidString() + fileExt
					saveFilePath := fmt.Sprintf("%s/%s", fPath, fName)
					vhost := viper.GetString("app.vhost")
					if strings.TrimSpace(vhost) == "" {
						vhost = "/"
					}
					head_img_url := fmt.Sprintf("%s%suploads/%s/%s", image_domain, vhost, time.Now().Format("2006-01"), fName)
					if pErr := os.MkdirAll(fPath, os.ModePerm); pErr == nil {
						if imgContent, rErr := ioutil.ReadAll(resp.Body); rErr == nil {
							if img, _, dErr := image.Decode(bytes.NewReader(imgContent)); dErr == nil {
								if dstImage := imaging.Resize(img, max_width, max_height, imaging.Linear); dstImage != nil {
									if sErr := imaging.Save(dstImage, saveFilePath, imaging.PNGCompressionLevel(png.BestSpeed)); sErr == nil {
										wxAuthInfo.HeadImgURL = head_img_url
									}
								}
							}
						}
					}
				}
			}

			returnURL := m.returnUrl
			if customReturnURL != "" {
				returnURL = customReturnURL
			}
			claims := pkg.WxAuthInfoClaims{
				WxAuthInfo: *wxAuthInfo,
			}
			m.jwtInstance.CreateToken(ctx, claims)
			parmSplitor := "?"
			if strings.Contains(returnURL, "?") {
				parmSplitor = "&"
			}
			url := ""
			if cak != "" {
				url = fmt.Sprintf("%s%scak=%s", returnURL, parmSplitor, cak)
			} else {
				url = returnURL
			}
			if afterAuth != nil {
				afterAuth(wxAuthInfo, cak)
			}
			ctx.Redirect(url)
		}
	}
}
