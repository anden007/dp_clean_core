package request

import (
	"encoding/json"

	"github.com/anden007/dp_clean_core/pkg/union_jd_sdk"
)

// UnionOpenPromotionBySubUnionidRequest 社交媒体获取推广链接
type UnionOpenOrderRowQueryRequest struct {
	OrderReq *OrderReq `json:"orderReq"`
}

func NewUnionOpenOrderRowQueryRequest(timeType int, startTime string, endTime string, pageIndex int, pageSize int) *UnionOpenOrderRowQueryRequest {
	req := NewOrderReq(timeType, startTime, endTime, pageIndex, pageSize)
	return &UnionOpenOrderRowQueryRequest{
		OrderReq: req,
	}
}

func (req *UnionOpenOrderRowQueryRequest) JsonParams() (string, error) {
	orderReq := map[string]interface{}{
		"orderReq": &req.OrderReq,
	}
	paramJsonBytes, err := json.Marshal(&orderReq)
	if err != nil {
		return "", err
	}
	return string(paramJsonBytes), nil
}

func (req *UnionOpenOrderRowQueryRequest) ResponseName() string {
	return "jd_union_open_order_row_query_responce"
}

func (req *UnionOpenOrderRowQueryRequest) GetResultFieldName() string {
	return "queryResult"
}

func (req *UnionOpenOrderRowQueryRequest) GetMethodName() string {
	return union_jd_sdk.MethodQueryBonus
}
