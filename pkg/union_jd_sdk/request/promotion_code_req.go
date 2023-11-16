package request

type PromotionCodeReq struct {
	MaterialId    string `json:"materialId"`              // [必填] 推广物料url，例如活动链接、商品链接、联盟链接（包含微信shortlink形式）等；不支持仅传入skuid
	SubUnionId    string `json:"subUnionId,omitempty"`    // 子渠道标识，仅支持传入字母、数字、下划线或中划线，最多80个字符（不可包含空格），该参数会在订单行查询接口中展示（需申请权限，申请方法请见https://union.jd.com/helpcenter/13246-13247-46301）
	PositionId    int64  `json:"positionId,omitempty"`    // 推广位ID
	Pid           string `json:"pid,omitempty"`           // 联盟子推客身份标识（不能传入接口调用者自己的pid）
	CouponUrl     string `json:"couponUrl,omitempty"`     // 优惠券领取链接，在使用优惠券、商品二合一功能时入参，且materialId须为商品详情页链接
	ChainType     int    `json:"chainType,omitempty"`     // 转链类型，1：长链， 2 ：短链 ，3： 长链+短链，默认短链，短链有效期60天
	GiftCouponKey string `json:"giftCouponKey,omitempty"` // 礼金批次号
	ChannelId     int64  `json:"channelId,omitempty"`     // 渠道关系ID
	Command       int    `json:"command,omitempty"`       // 是否生成短口令，1：生成，默认不生成（需申请权限，申请方法请见https://union.jd.com/helpcenter/13246-13247-46301）
	WeChatType    int    `json:"weChatType,omitempty"`    // 微信小程序ShortLink类型（需向cps-qxsq@jd.com申请权限）
}

func NewPromotionCodeReq(materialId string, subUnionId string, positionId int64, couponUrl string, chainType int, channelId int64) *PromotionCodeReq {
	if chainType == 0 {
		chainType = 1 // 默认长链
	}
	return &PromotionCodeReq{
		MaterialId:    materialId,
		SubUnionId:    subUnionId,
		PositionId:    positionId,
		Pid:           "",
		CouponUrl:     couponUrl,
		ChainType:     chainType,
		GiftCouponKey: "",
		ChannelId:     channelId,
		Command:       0,
		WeChatType:    0,
	}
}
