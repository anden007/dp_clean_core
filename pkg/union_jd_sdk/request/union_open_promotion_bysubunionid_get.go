package request

import (
	"encoding/json"

	"github.com/anden007/af_dp_clean_core/pkg/union_jd_sdk"
)

// UnionOpenPromotionBySubUnionidRequest 社交媒体获取推广链接
type UnionOpenPromotionBySubUnionidRequest struct {
	PromotionCodeReq *PromotionCodeReq `json:"promotionCodeReq"`
}

func NewUnionOpenPromotionBySubUnionidRequest(materialId string, subUnionId string, positionId int64, couponUrl string, chainType int, channelId int64) *UnionOpenPromotionBySubUnionidRequest {
	req := NewPromotionCodeReq(materialId, subUnionId, positionId, couponUrl, chainType, channelId)
	return &UnionOpenPromotionBySubUnionidRequest{
		PromotionCodeReq: req,
	}
}

func (req *UnionOpenPromotionBySubUnionidRequest) JsonParams() (string, error) {
	promotionCodeReq := map[string]interface{}{
		"promotionCodeReq": &req.PromotionCodeReq,
	}
	paramJsonBytes, err := json.Marshal(&promotionCodeReq)
	if err != nil {
		return "", err
	}
	return string(paramJsonBytes), nil
}

func (req *UnionOpenPromotionBySubUnionidRequest) ResponseName() string {
	return "jd_union_open_promotion_bysubunionid_get_responce"
}

func (req *UnionOpenPromotionBySubUnionidRequest) GetResultFieldName() string {
	return "getResult"
}

func (req *UnionOpenPromotionBySubUnionidRequest) GetMethodName() string {
	return union_jd_sdk.MethodBysubunionidPromotion
}
