package response

type PromotionCodeResp struct {
	ShortURL        string `json:"shortUrl"`        // 生成的推广目标链接，以短链接形式，有效期60天
	ClickURL        string `json:"clickUrl"`        // 生成推广目标的长链，长期有效
	JCommand        string `json:"jCommand"`        // 需要权限申请，京口令（匹配到红包活动有效配置才会返回京口令）
	JShortCommand   string `json:"jShortCommand"`   // 需要权限申请，短口令
	WeChatShortLink string `json:"weChatShortLink"` // 微信小程序ShortLink（需向cps-qxsq@jd.com申请权限）
}

type PromotionCodeResult struct {
	UnionJDApiResult
	Data *PromotionCodeResp `json:"data"`
}
