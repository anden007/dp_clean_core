package part

import (
	"fmt"
	"net/url"
	"time"

	"github.com/anden007/af_dp_clean_core/pkg/union_jd_sdk"
	"github.com/anden007/af_dp_clean_core/pkg/union_jd_sdk/common"
	"github.com/anden007/af_dp_clean_core/pkg/union_jd_sdk/response"

	jsoniter "github.com/json-iterator/go"
	jsonTime "github.com/liamylian/jsontime/v2/v2"
	"github.com/spf13/viper"
)

type IUnionJD interface {
	BysubunionidPromotion(req common.Request) (result response.PromotionCodeResult, err error)
	OrderRowQuery(req common.Request) (result response.OrderRowQueryResult, err error)
	GetJDDeepLink(platform string, sourceUrl string) (result string, err error)
}

type UnionJD struct {
	enable      bool
	accessToken string
	appKey      string
	appSecret   string
	sdkClient   *union_jd_sdk.JdClient
	jsonEncoder jsoniter.API
}

func NewUnionJD() IUnionJD {
	// loadTime := time.Now()
	jsonEncoder := jsonTime.ConfigWithCustomTimeFormat
	jsonTime.SetDefaultTimeFormat("2006-01-02 15:04:05", time.Local)
	instance := &UnionJD{
		enable:      viper.GetBool("union-jd.enable"),
		accessToken: viper.GetString("union-jd.access_token"),
		appKey:      viper.GetString("union-jd.app_key"),
		appSecret:   viper.GetString("union-jd.app_secret"),
		jsonEncoder: jsonEncoder,
	}
	if instance.enable {
		instance.sdkClient = union_jd_sdk.NewJdClient(instance.accessToken, instance.appKey, instance.appSecret)
	}
	// if ENV == ENUM_ENV_DEV {
	// 	misc.ServiceLoadInfo("UnionJD", instance.enable, loadTime)
	// }
	return instance
}

// 社交媒体获取推广链接接口
func (m *UnionJD) BysubunionidPromotion(req common.Request) (result response.PromotionCodeResult, err error) {
	if m.enable {
		if resultByte, eErr := m.sdkClient.Execute(req); eErr == nil {
			if jErr := m.jsonEncoder.Unmarshal(resultByte, &result); jErr != nil {
				err = jErr
			}
		} else {
			err = eErr
		}
	} else {
		panic("UnionJD is not enabled")
	}
	return
}

// 查询推广订单及佣金信息
func (m *UnionJD) OrderRowQuery(req common.Request) (result response.OrderRowQueryResult, err error) {
	if m.enable {
		if resultByte, eErr := m.sdkClient.Execute(req); eErr == nil {
			if jErr := m.jsonEncoder.Unmarshal(resultByte, &result); jErr != nil {
				err = jErr
			}
		} else {
			err = eErr
		}
	} else {
		panic("UnionJD is not enabled")
	}
	return
}

// 生成京东Deeplink
func (m *UnionJD) GetJDDeepLink(platform string, sourceUrl string) (result string, err error) {
	if sourceUrl != "" {
		resultUrl, _ := url.Parse(sourceUrl)
		args := resultUrl.Query()
		args.Set("e", fmt.Sprintf("af^%s^__ACCOUNTID__^__DID__^__AID__^__CID__^__CALLBACK_PARAM__", platform))
		resultUrl.RawQuery = args.Encode()
		paramsStr := url.QueryEscape(fmt.Sprintf("{\"category\":\"jump\",\"des\":\"m\",\"url\":\"%s\",\"keplerID\":\"kpl_jdjdmyl00000001\"}", resultUrl.String()))
		result = fmt.Sprintf("openapp.jdmobile://virtual?params=%s", paramsStr)
	}
	return
}
