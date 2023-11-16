package union_jd_sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/anden007/dp_clean_core/pkg/union_jd_sdk/common"
)

const ServerUrl = "https://api.jd.com/routerjson"

type JdClient struct {
	accessToken string
	appKey      string
	appSecret   string
}

func NewJdClient(accessToken, appKey, appSecret string) *JdClient {
	return &JdClient{accessToken: accessToken, appKey: appKey, appSecret: appSecret}
}

func (c *JdClient) Execute(req common.Request) ([]byte, error) {
	var result []byte
	method := req.GetMethodName()
	// get business params
	jsonParams, err := req.JsonParams()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[UnionJD] 获取JsonParams失败, 错误:%s", err.Error()))
	}
	// part.Log().Debug("[UnionJD] 业务参数", zap.String("JsonParams", jsonParams))

	// sign
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	signParams := map[string]string{
		"app_key":           c.appKey,
		"format":            "json",
		"method":            method,
		"360buy_param_json": jsonParams,
		"timestamp":         timestamp,
		"v":                 "1.0",
	}
	signValue := common.Sign(signParams, c.appSecret)
	params := common.Config{
		Version:     "1.0",
		Method:      method,
		AccessToken: c.accessToken,
		AppKey:      c.appKey,
		Format:      "json",
		Timestamp:   timestamp,
		Sign:        signValue,
		ParamJson:   jsonParams,
	}
	// part.Log().Debug("[UnionJD] 请求参数", zap.Any("data", params))

	// 请求JD Api服务器
	respBytes, err := common.HttpGet(ServerUrl, params)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[UnionJD] HTTP请求失败, 错误:%s", err.Error()))
	}

	var respObj map[string]interface{}
	if err := json.Unmarshal(respBytes, &respObj); err != nil {
		return nil, errors.New(fmt.Sprintf("[UnionJD] JSON反序列化失败, 错误:%s", err.Error()))
	}

	if _, exists := respObj[req.ResponseName()]; exists {
		responseMessage := respObj[req.ResponseName()].(map[string]interface{})
		// respCode := responseMessage["code"].(string)
		respResult := responseMessage[req.GetResultFieldName()].(string)
		result = []byte(respResult)
		// part.Log().Debug("[UnionJD] 响应结果", zap.String("code", respCode))
	} else {
		result = respBytes
	}
	return result, nil
}
