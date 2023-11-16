package common

type Request interface {
	JsonParams() (string, error)
	ResponseName() string
	GetMethodName() string
	GetResultFieldName() string
}

type Result struct {
	Code       int32       `json:"code"`
	Message    string      `json:"message"`
	RequestId  string      `json:"requestId"`
	TotalCount int64       `json:"totalCount"`
	Data       interface{} `json:"data"`
}
