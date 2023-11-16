package request

import (
	"encoding/json"

	"github.com/anden007/dp_clean_core/pkg/union_jd_sdk"
)

// UnionOpenGoodsJingfenQuery 京粉精选商品查询接口
type UnionOpenGoodsJingfenQueryRequest struct {
	GoodsReq *JFGoodsReq `json:"goods_req"`
}

func NewUnionOpenGoodsJingfenQueryRequest(goodsReq *JFGoodsReq) *UnionOpenGoodsJingfenQueryRequest {
	return &UnionOpenGoodsJingfenQueryRequest{
		GoodsReq: goodsReq,
	}
}

func (req *UnionOpenGoodsJingfenQueryRequest) JsonParams() (string, error) {
	goodsReq := map[string]interface{}{
		"goodsReq": &req.GoodsReq,
	}
	paramJsonBytes, err := json.Marshal(&goodsReq)
	if err != nil {
		return "", err
	}
	return string(paramJsonBytes), nil
}

func (req *UnionOpenGoodsJingfenQueryRequest) ResponseName() string {
	return "jd_union_open_goods_jingfen_query_response"
}

func (req *UnionOpenGoodsJingfenQueryRequest) GetResultFieldName() string {
	return "queryResult"
}

func (req *UnionOpenGoodsJingfenQueryRequest) GetMethodName() string {
	return union_jd_sdk.MethodQueryJingfenGoods
}
